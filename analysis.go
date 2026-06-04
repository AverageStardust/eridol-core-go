package main

import (
	"encoding/binary"
	"math"
	"sync"

	"github.com/madelynnblue/go-dsp/fft"
)

const octaves = 6

const fftSize = 355

const noise1Bin = 27
const bBin = 29
const noise2Bin = 32
const noise3Bin = 33
const dsBin = 36
const noise4Bin = 39
const noise5Bin = 41
const fsBin = 44
const claimBin = 49
const aBin = 52
const counterClaimBin = 55
const noise6Bin = 57

var inputBuffer *ring[float64]
var inputMutex *sync.Mutex = &sync.Mutex{}
var octaveBuffers [octaves]*ring[float64]
var octaveSounds [octaves]OctaveSound

func init() {
	inputBuffer = newRing[float64](1 << 13)

	for i := range octaves {
		octaveBuffers[i] = newRing[float64](1 << 13)
		octaveSounds[i] = OctaveSound{}
	}
}

func enqueueData(inBuffer []byte, frameCount uint32) {
	inputMutex.Lock()
	for i := range frameCount {
		bits := binary.NativeEndian.Uint16(inBuffer[i*2:])
		amp := float64(int16(bits)) / math.MaxInt16
		inputBuffer.Enqueue(amp)
	}
	inputMutex.Unlock()
}

func analyze() {
	sampleDate()

analysisLoop:
	for {
		for octave := octaves - 1; octave >= 0; octave-- {
			// dequeue data
			timeDomain, success := octaveBuffers[octave].DequeueBatch(fftSize)

			if !success {
				if octave == 5 {
					break analysisLoop
				} else {
					break
				}
			}

			// run fft
			freqDomain := fft.FFTReal(timeDomain)

			// summerize data
			sound := newOctaveSound(freqDomain)

			if octave == 0 {
				// don't average to make lowest octave more responsive (and less acurate)
				octaveSounds[octave] = sound
			} else {
				// take average to make more acurate
				octaveSounds[octave] = octaveSounds[octave].add(sound).scale(0.5)
			}
		}
	}
}

// get enough data from the ring buffer to do FFTs for each octave
func sampleDate() {
	// take an amount of data that can be downsampled evenly, discard the extra
	usableSampleCount := int(inputBuffer.Size()) & (-1 << (octaves - 1))

	inputMutex.Lock()
	usableSampels, _ := inputBuffer.DequeueBatch(uint64(usableSampleCount))
	inputMutex.Unlock()

	samplings := make([][]float64, octaves)
	samplings[octaves-1] = usableSampels

	// create different down-samplings at an idea sampling rate for each octave
	for i := octaves - 1; i > 0; i-- {
		samplings[i-1] = halfDownSample(samplings[i])
	}

	for i := range octaves {
		octaveBuffers[i].EnqueueBatch(samplings[i])
	}
}

// take sound data and make it have half the sampling rate
func halfDownSample(input []float64) []float64 {
	var output = make([]float64, len(input)/2)

	for i := range len(input) / 2 {
		output[i] = (input[i*2] + input[i*2+1]) / 2
	}

	return output
}

func (a OctaveSound) add(b OctaveSound) OctaveSound {
	return OctaveSound{
		B:            a.B + b.B,
		Ds:           a.Ds + b.Ds,
		Fs:           a.Fs + b.Fs,
		A:            a.A + b.A,
		Claim:        a.Claim + b.Claim,
		CounterClaim: a.CounterClaim + b.CounterClaim,
		Noise:        a.Noise + b.Noise,
	}
}

func (a OctaveSound) scale(b float64) OctaveSound {
	return OctaveSound{
		B:            a.B * b,
		Ds:           a.Ds * b,
		Fs:           a.Fs * b,
		A:            a.A * b,
		Claim:        a.Claim * b,
		CounterClaim: a.CounterClaim * b,
		Noise:        a.Noise * b,
	}
}
