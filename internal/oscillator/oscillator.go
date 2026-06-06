package oscillator

import (
	"encoding/binary"
	"math"
)

type plan struct {
}

type oscillator struct {
	octaveSounds []octaveSound
	timeSamples  uint64
}

func newOscillator(octaves int) *oscillator {
	return &oscillator{
		octaveSounds: make([]octaveSound, octaves),
	}
}

func (osc oscillator) Write(buffer []byte, samples, sampleRate int) {
	for i := range samples {
		time := float64(osc.timeSamples) / float64(sampleRate)

		var amplitude float32
		for octave, sound := range osc.octaveSounds {
			amplitude += sound.Sample(time, octave)
		}

		amplitude /= osc.totalAmplitude()

		bits := uint16(int16(amplitude * math.MaxInt16))
		binary.NativeEndian.PutUint16(buffer[i*2:], bits)

		osc.timeSamples++
	}
}

func (osc oscillator) totalAmplitude() float32 {
	var total float32
	for _, sound := range osc.octaveSounds {
		total += sound.TotalAmplitude()
	}

	return total
}
