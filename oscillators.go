package core

import (
	"encoding/binary"
	"log"
	"math"
)

var oscillatorSamples uint64
var afterLastOscillatorGoesSilent func() = nil

func writeOscillators(buffer []byte, frameCount int) {
	oscillatorMutex.Lock()
	defer oscillatorMutex.Unlock()

	for i := range frameCount {
		amplitude := sumAndNormalizeOscillators()

		bits := uint16(int16(amplitude * math.MaxInt16))
		binary.NativeEndian.PutUint16(buffer[i*2:i*2+2], bits)

		oscillatorSamples++
	}

	if DoLogging {
		log.Println("eridol-core: Wrote ", frameCount, " frames of audio input")
	}
}

func sumAndNormalizeOscillators() float32 {
	var totalAmplitude, totalMax float32

	for oscillator := range oscillators {
		amplitude, max := oscillator.sample()

		totalAmplitude += amplitude
		totalMax += max
	}

	if totalMax < 1 {
		totalMax = 1
	}

	return totalAmplitude / totalMax
}

func closeOscillators(callback func()) {
	oscillatorMutex.Lock()
	defer oscillatorMutex.Unlock()

	for osc := range oscillators {
		osc.Close()
	}

	afterLastOscillatorGoesSilent = callback
}

func isClosingOscillators() bool {
	return afterLastOscillatorGoesSilent != nil
}
