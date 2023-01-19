package snd

import (
	"math"
)

type Tone struct {
	Start             uint16    `json:"start,omitempty"`
	Length            uint16    `json:"length,omitempty"` // default=500
	FreqBase          *Envelope `json:"freq_base,omitempty"`
	FreqModRate       *Envelope `json:"freq_mod_rate,omitempty"`
	FreqModRange      *Envelope `json:"freq_mod_range,omitempty"`
	AmpBase           *Envelope `json:"amp_base,omitempty"`
	AmpModRate        *Envelope `json:"amp_mod_rate,omitempty"`
	AmpModRange       *Envelope `json:"amp_mod_range,omitempty"`
	HarmonicVolumes   [5]int    `json:"harmonic_volumes,omitempty"`
	HarmonicSemitones [5]int    `json:"harmonic_semitones,omitempty"`
	HarmonicDelays    [5]int    `json:"harmonic_delays,omitempty"`
	Filter            *filter   `json:"filter,omitempty"`
	FilterRange       *Envelope `json:"filter_range,omitempty"`
	Release           *Envelope `json:"release,omitempty"`
	Attack            *Envelope `json:"attack,omitempty"`
	ReverbDelay       int       `json:"reverb_delay,omitempty"`
	ReverbVolume      int       `json:"reverb_volume,omitempty"`
}

func NewTone() *Tone {
	t := &Tone{}
	t.Length = 500
	t.ReverbVolume = 100
	return t
}

func (t *Tone) Synthesize(sampleCount, duration int) (samples []int, err error) {
	if duration < 10 {
		return
	}

	samplesPerStep := float64(sampleCount) / float64(duration)

	t.FreqBase.reset()
	t.AmpBase.reset()

	var freqStart int
	var freqDuration int
	var freqPhase int

	if t.FreqModRate != nil {
		t.FreqModRate.reset()
		t.FreqModRange.reset()

		freqDuration = int((32.768 * float64(t.FreqModRate.End-t.FreqModRate.Start)) / samplesPerStep)
		freqStart = int((32.768 * float64(t.FreqModRate.Start)) / samplesPerStep)
	}

	var ampStart int
	var ampDuration int
	var ampPhase int

	if t.AmpModRate != nil {
		t.AmpModRate.reset()
		t.AmpModRange.reset()

		ampDuration = int((32.768 * float64(t.AmpModRate.End-t.AmpModRate.Start)) / samplesPerStep)
		ampStart = int((32.768 * float64(t.AmpModRate.Start)) / samplesPerStep)
	}

	var _phases [5]int
	var _delays [5]int
	var _volumes [5]int
	var _semitones [5]int
	var _starts [5]int

	for h := 0; h < 5; h++ {
		if t.HarmonicVolumes[h] != 0 {
			_phases[h] = 0
			_delays[h] = int(float64(t.HarmonicDelays[h]) * samplesPerStep)
			_volumes[h] = (t.HarmonicVolumes[h] << 14) / 100

			semitone := float64(t.FreqBase.End - t.FreqBase.Start)
			semitone *= 32.768 * math.Pow(1.0057929410678534, float64(t.HarmonicSemitones[h]))
			semitone /= samplesPerStep

			_semitones[h] = int(semitone)
			_starts[h] = int((32.768 * float64(t.FreqBase.Start)) / samplesPerStep)
		}
	}

	samples = make([]int, sampleCount)

	for position := 0; position < sampleCount; position++ {
		freq := t.FreqBase.eval(sampleCount)
		amp := t.AmpBase.eval(sampleCount)

		if t.FreqModRate != nil {
			rate := t.FreqModRate.eval(sampleCount)
			_range := t.FreqModRange.eval(sampleCount)
			freq += generate(t.FreqModRate.Form, freqPhase, _range) >> 1
			freqPhase += ((rate * freqDuration) >> 16) + freqStart
		}

		if t.AmpModRate != nil {
			rate := t.AmpModRate.eval(sampleCount)
			_range := t.AmpModRange.eval(sampleCount)
			amp = (amp * ((generate(t.AmpModRate.Form, ampPhase, _range) >> 1) + 32768)) >> 15
			ampPhase += ((rate * ampDuration) >> 16) + ampStart
		}

		for h := 0; h < 5; h++ {
			if t.HarmonicVolumes[h] != 0 {
				offset := position + _delays[h]

				if offset >= sampleCount {
					continue
				}

				samples[offset] += generate(t.FreqBase.Form, _phases[h], (amp*_volumes[h])>>15)
				_phases[h] += ((freq * _semitones[h]) >> 16) + _starts[h]
			}
		}
	}

	if t.Release != nil {
		t.Release.reset()
		t.Attack.reset()

		counter := 0
		muted := true

		for position := 0; position < sampleCount; position++ {
			release := t.Release.eval(sampleCount)
			attack := t.Attack.eval(sampleCount)
			threshold := 0

			if muted {
				threshold = t.Release.Start + (((t.Release.End - t.Release.Start) * release) >> 8)
			} else {
				threshold = t.Release.Start + (((t.Release.End - t.Release.Start) * attack) >> 8)
			}

			counter += 256

			if counter >= threshold {
				counter = 0
				muted = !muted
			}

			if muted {
				samples[position] = 0
			}
		}
	}

	if t.ReverbDelay > 0 && t.ReverbVolume > 0 {
		start := int(float64(t.ReverbDelay) * samplesPerStep)

		for position := start; position < sampleCount; position++ {
			samples[position] += (samples[position-start] * t.ReverbVolume) / 100
		}
	}

	if t.Filter.Poles[0] > 0 || t.Filter.Poles[1] > 0 {
		t.FilterRange.reset()
		_range := t.FilterRange.eval(sampleCount + 1)
		fwd := t.Filter.eval(0, float64(_range)/65536.0)
		aft := t.Filter.eval(1, float64(_range)/65536.0)

		if sampleCount >= fwd+aft {
			index := 0
			interval := aft

			if aft > sampleCount-fwd {
				interval = sampleCount - fwd
			}

			for index < interval {
				sample := int((int64(samples[index+fwd]) * _unity16) >> 16)

				for offset := 0; offset < fwd; offset++ {
					sample += (int)((int64(samples[index+fwd-1-offset]) * _coef16[0][offset]) >> 16)
				}

				for offset := 0; offset < index; offset++ {
					sample -= int((int64(samples[index-1-offset]) * _coef16[1][offset]) >> 16)
				}

				samples[index] = sample
				_range = t.FilterRange.eval(sampleCount + 1)
				index++
			}

			interval = 128

			for {
				if interval > sampleCount-fwd {
					interval = sampleCount - fwd
				}

				for index < interval {
					sample := int((int64(samples[index+fwd]) * _unity16) >> 16)

					for offset := 0; offset < fwd; offset++ {
						sample += (int)((int64(samples[index+fwd-1-offset]) * _coef16[0][offset]) >> 16)
					}

					for offset := 0; offset < aft; offset++ {
						sample -= (int)((int64(samples[index-1-offset]) * _coef16[1][offset]) >> 16)
					}

					samples[index] = sample
					_range = t.FilterRange.eval(sampleCount + 1)
					index++
				}

				if index >= sampleCount-fwd {
					for index < sampleCount {
						sample := 0

						for offset := index + fwd - sampleCount; offset < fwd; offset++ {
							sample += (int)((int64(samples[index+fwd-1-offset]) * _coef16[0][offset]) >> 16)
						}

						for offset := 0; offset < aft; offset++ {
							sample -= (int)((int64(samples[index-1-offset]) * _coef16[1][offset]) >> 16)
						}

						samples[index] = sample
						t.FilterRange.eval(sampleCount + 1)
						index++
					}
					break
				}

				fwd = t.Filter.eval(0, float64(_range)/65536.0)
				aft = t.Filter.eval(1, float64(_range)/65536.0)
				interval += 128
			}
		}
	}

	for i, sample := range samples {
		if sample > 32767 {
			samples[i] = 32767
		} else if sample < -32768 {
			samples[i] = -32768
		}
	}

	return
}

func (t *Tone) read(in *buffer) error {
	t.FreqBase = &Envelope{}
	t.AmpBase = &Envelope{}

	if err := t.FreqBase.read(in); err != nil {
		return err
	} else if err := t.AmpBase.read(in); err != nil {
		return err
	}

	if in.u8() != 0 {
		in.rewind(1)
		t.FreqModRate = &Envelope{}
		t.FreqModRange = &Envelope{}

		if err := t.FreqModRate.read(in); err != nil {
			return err
		} else if err := t.FreqModRange.read(in); err != nil {
			return err
		}
	}

	if in.u8() != 0 {
		in.rewind(1)
		t.AmpModRate = &Envelope{}
		t.AmpModRange = &Envelope{}

		if err := t.AmpModRate.read(in); err != nil {
			return err
		} else if err := t.AmpModRange.read(in); err != nil {
			return err
		}
	}

	if in.u8() != 0 {
		in.rewind(1)
		t.Release = &Envelope{}
		t.Attack = &Envelope{}

		if err := t.Release.read(in); err != nil {
			return err
		} else if err := t.Attack.read(in); err != nil {
			return err
		}
	}

	// harmonics
	for harmony := 0; harmony < 10; harmony++ {
		volume := in.usmart()

		if volume == 0 {
			break
		}

		t.HarmonicVolumes[harmony] = volume
		t.HarmonicSemitones[harmony] = in.smart()
		t.HarmonicDelays[harmony] = in.usmart()
	}

	t.ReverbDelay = in.usmart()
	t.ReverbVolume = in.usmart()
	t.Length = in.u16()
	t.Start = in.u16()

	t.Filter = &filter{}
	t.FilterRange = &Envelope{}
	return t.Filter.read(in, t.FilterRange)
}
