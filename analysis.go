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
	sampleData()

analysisLoop:
	for {
		for octave := octaves - 1; octave >= 0; octave-- {
			// dequeue data
			timeDomain, success := octaveBuffers[octave].DequeueBatch(fftSize)

			if !success {
				if octave == octaves-1 {
					break analysisLoop
				} else {
					break
				}
			}

			// run fft
			freqDomain := applyFFT(timeDomain)

			// summerize data
			sound := newOctaveSound(freqDomain)

			if octave == 0 {
				// don't average to make lowest octave more responsive (and less acurate)
				octaveSounds[octave] = sound
			} else {
				// take average to make more acurate
				octaveSounds[octave] = octaveSounds[octave].Add(sound).Scale(0.5)
			}
		}
	}
}

// run fft on
func applyFFT(timeDomain []float64) []complex128 {
	return fft.FFTReal(timeDomain)
}

// downsample data from the input buffer into a buffer ideal for each octave
func sampleData() {
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
