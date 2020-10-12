package snd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

const SampleRate = 22050

func LoadTracks(reader io.Reader) (tracks map[int]*Track, err error) {
	data, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	tracks = make(map[int]*Track)
	buf := &buffer{data, 0}
	count := 0

	for {
		index := int(buf.u16())

		if index == 65535 {
			break
		}

		if track, err := loadTrack(buf); err != nil {
			return nil, fmt.Errorf("%s: sound[%d] (counter=%d)", err.Error(), index, count)
		} else {
			tracks[index] = track
		}

		count++
	}

	return
}

func LoadTrack(filename string) (*Track, error) {
	if data, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else if track, err := loadTrack(&buffer{data, 0}); err != nil {
		return nil, err
	} else {
		return track, nil
	}
}

func loadTrack(b *buffer) (*Track, error) {
	t := &Track{}
	if err := t.read(b); err != nil {
		return nil, err
	}
	return t, nil
}

type Track struct {
	Tones     [10]*tone
	Delay     int
	LoopBegin uint16
	LoopEnd   uint16
}

func (t *Track) read(in *buffer) error {
	for i := 0; i < 10; i++ {
		if in.u8() != 0 {
			in.rewind(1)
			t.Tones[i] = newTone()
			if err := t.Tones[i].read(in); err != nil {
				return err
			}
		}
	}

	t.LoopBegin = in.u16()
	t.LoopEnd = in.u16()
	return nil
}

func (t *Track) generate() ([]byte, error) {
	length := 0

	for i := 0; i < len(t.Tones); i++ {
		if t := t.Tones[i]; t != nil {
			length = int(math.Max(float64(t.Start+t.Length), float64(length)))
		}
	}

	if length == 0 {
		return nil, errors.New("empty sound")
	}

	sampleCount := (length * SampleRate) / 1000
	samples := make([]int, sampleCount)

	for _, tone := range t.Tones {
		if tone != nil {
			toneSampleCount := (int(tone.Length) * SampleRate) / 1000
			toneStart := (int(tone.Start) * SampleRate) / 1000
			toneSamples, err := tone.generate(toneSampleCount, int(tone.Length))

			if err != nil {
				return nil, err
			}

			for pos := 0; pos < toneSampleCount; pos++ {
				samples[pos+toneStart] += toneSamples[pos]
			}
		}
	}

	sampleCount *= 2

	buf := &buffer{make([]byte, 44+sampleCount), 0}
	buf.put([]byte("RIFF"))
	buf.p32le(36 + sampleCount)
	buf.put([]byte("WAVE"))
	buf.put([]byte("fmt "))
	buf.p32le(16)    // Subchunk 1 Size
	buf.p16le(1)     // PCM Format
	buf.p16le(1)     // Mono
	buf.p32le(22050) // Sample Rate
	buf.p32le(44100) // Byte Rate
	buf.p16le(2)     // BlockAlign
	buf.p16le(16)    // BitsPerSample
	buf.put([]byte("data"))
	buf.p32le(sampleCount)

	for _, sample := range samples {
		buf.p16le(int(sample))
	}

	return buf.data, nil
}
