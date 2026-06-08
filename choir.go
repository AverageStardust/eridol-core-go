package core

import (
	"encoding/binary"
	"math"
)

// A set of synthesizers, one for each octave.
type Choir struct {
	synths      [OctaveCount]*Synthesizer
	timeSamples uint64
}

var choir = newChoir()

func newChoir() *Choir {
	choir := &Choir{}

	for i := range OctaveCount {
		choir.synths[i] = newSynth(i)
	}

	return choir
}

// Returns the synthesizer for the octave Nth octave.
// This synth is the same for the entire life of the program.
func Synth(n int) *Synthesizer {
	return choir.synths[n]
}

func (choir *Choir) writeTo(quanta soundQuanta) {
	for i := range quanta.frameCount {
		amplitude := choir.sample()

		bits := uint16(int16(amplitude * math.MaxInt16))
		binary.NativeEndian.PutUint16(quanta.buffer[i*2:i*2+2], bits)
	}
}

func (choir *Choir) sample() float32 {
	var totalAmplitude, totalMax float32

	for _, synth := range choir.synths {
		synth.mutex.Lock()
	}

	for _, synth := range choir.synths {
		amplitude, max := synth.sample(choir.timeSamples)

		totalAmplitude += amplitude
		totalMax += max
	}

	for _, synth := range choir.synths {
		synth.mutex.Unlock()
	}

	choir.timeSamples++

	if totalMax < 1 {
		totalMax = 1
	}

	return totalAmplitude / totalMax
}

func (choir *Choir) silence() {
	for _, synth := range choir.synths {
		synth.Silence()
	}
}
