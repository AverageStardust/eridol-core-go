package core

import (
	"encoding/binary"
	"math"
)

type oscillator struct {
	samplesPerSecond int
	sampleTime       int64
}

func newOscillator(samplesPerSecond int) *oscillator {
	return &oscillator{
		samplesPerSecond: samplesPerSecond,
		sampleTime:       0,
	}
}

func (osc *oscillator) write(buffer []byte, samples int) {
	for i := range samples {
		freq := 440.0
		phase := float64(osc.sampleTime) / float64(osc.samplesPerSecond) * freq
		bits := math.Float32bits(float32(math.Sin(phase * math.Pi * 2)))
		binary.NativeEndian.PutUint32(buffer[i*4:], bits)
		osc.sampleTime++
	}
}
