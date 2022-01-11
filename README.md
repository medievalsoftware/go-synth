# Simple Waveform Synthesizer

# Format:

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

	Release				Envelope
	Attack				Envelope

	Harmonics			[10]Harmonic

	ReverbDelay			usmart
	ReverbVolume		usmart
	Length				uint16
	Start				uint16

	// Infinite Impulse Response Filter (IIR Filter)
	IIRFilter			Filter
}

type Envelope struct {
	Form		uint8 // 1 = Square, 2 = Sine, 3 = Saw, 4 = Noise
	// if Form == 0 then skip remaining:
	Start		uint32
	End			uint32
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
	Frequencies [2][2][4]uint16
	Ranges      [2][2][4]uint16
	Ranges		Envelope
}
```

# Special types:
```go
// range: [-16384, 16383]
func (b *buffer) smart() int {
	if b.data[b.position] < 128 {
		return int(b.u8()) - 64
	}
	return int(b.u16()) - 49152
}

// range: [0, 32767]
func (b *buffer) usmart() int {
	if b.data[b.position] < 128 {
		return int(b.u8())
	}
	return int(b.u16()) - 32768
}
```
