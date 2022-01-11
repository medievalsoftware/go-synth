# Simple Waveform Synthesizer

# Description

This is a 16bit mono PCM WAVE generator that supports looping, modulation(of frequency and amplitude), release, attack, harmonics, and reverb.

# Input Format:

```go
type Track struct {
	Tones		[10]Tone
	LoopBegin	uint16
	LoopEnd		uint16
}

type Tone struct {
	Active				bool // uint8

	// if !Active then skip remaining:
	FrequencyBase		Envelope
	AmplitudeBase		Envelope

	FrequencyModRate	Envelope
	FrequencyModRange	Envelope

	AmplitudeModRate	Envelope
	AmplitudeModRange	Envelope

	Release			Envelope
	Attack			Envelope

	Harmonics		[10]Harmonic

	ReverbDelay		usmart
	ReverbVolume		usmart
	Length			uint16
	Start			uint16

	// Infinite Impulse Response Filter (IIR Filter)
	IIRFilter		Filter
}

type Envelope struct {
	Form		uint8 // 1 = Square, 2 = Sine, 3 = Saw, 4 = Noise
	// if Form == 0 then skip remaining:
	
	Start		uint32
	End		uint32
	SegmentN	uint8
	Segments	[SegmentN]EnvelopeSegment
}

type EnvelopeSegment struct {
	Duration	uint16
	Peak		uint16
}

type Harmonic struct {
	Volume		usmart
	// if Volume == 0 then skip remaining:
	
	Semitone	smart
	Delay		usmart
}

// Infinite Impulse Response Filter (IIR Filter)
type Filter struct {
	Interval	uint8
	// if Interval == 0 then skip remaining:
	
	Unities		[2]uint16
	Frequencies	[2][2][4]uint16
	Ranges		[2][2][4]uint16
	Ranges		Envelope
}
```

## Special types:
```go
type smart	//(int8 or int16)	range: [-16384,	16383]
type usmart	//(uint8 or uint16)	range: [0,	32767]

func (b *buffer) smart() int {
	if b.data[b.position] < 128 {
		return int(b.u8()) - 64
	}
	return int(b.u16()) - 49152
}

func (b *buffer) usmart() int {
	if b.data[b.position] < 128 {
		return int(b.u8())
	}
	return int(b.u16()) - 32768
}
```

# FAQ
## What format does it output?
```
Format:         PCM
Channels:       1 (Mono)
Sample Rate:    22050Hz
Byte Rate:      44100Hz
BlockAlign:     2
BitsPerSample:  16
```

http://soundfile.sapp.org/doc/WaveFormat/

## What is this from?
Reverse engineered from a RuneScape Game Client from 2005.

## How do I use it?
At the moment I haven't written a tool to convert a human readable format (JSON) into the correct format. This means you would have to manually build the binary format yourself or create a tool to do it.

