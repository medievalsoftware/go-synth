package snd

import "fmt"

type envelope struct {
	Form      uint8
	Start     int
	End       int
	Length    int
	Durations []int
	Peaks     []int

	threshold int
	delta     int
	amplitude int
	ticks     int
	position  int
}

func (e *envelope) read(in *buffer) error {
	e.Form = in.u8()

	if e.Form > 5 {
		return fmt.Errorf("invalid envelope form: %d", e.Form)
	}

	e.Start = int(in.u32())
	e.End = int(in.u32())

	return e.readShape(in)
}

func (e *envelope) readShape(in *buffer) error {
	length := in.u8()

	if length == 0 {
		return fmt.Errorf("envelope with no shape")
	}

	e.Length = int(length)
	e.Durations = make([]int, length)
	e.Peaks = make([]int, length)

	for i := uint8(0); i < length; i++ {
		e.Durations[i] = int(in.u16())
		e.Peaks[i] = int(in.u16())
	}
	return nil
}

func (e *envelope) eval(delta int) int {
	// no shape
	if len(e.Peaks) == 0 {
		return 0
	}

	if e.ticks >= e.threshold {
		e.amplitude = e.Peaks[e.position] << 15
		e.position++

		if e.position >= e.Length {
			e.position = e.Length - 1
		}

		e.threshold = int((float64(e.Durations[e.position]) / 65536.0) * float64(delta))

		if e.threshold > e.ticks {
			e.delta = ((e.Peaks[e.position] << 15) - e.amplitude) / (e.threshold - e.ticks)
		}
	}

	e.amplitude += e.delta
	e.ticks++
	return (e.amplitude - e.delta) >> 15
}

func (e *envelope) reset() {
	e.threshold = 0
	e.position = 0
	e.delta = 0
	e.amplitude = 0
	e.ticks = 0
}
