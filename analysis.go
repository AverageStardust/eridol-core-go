package main

import (
	"encoding/binary"
	"log"
	"math"
	"sync"

	algofft "github.com/cwbudde/algo-fft"
)

const octaves = 6

const fftSize = 374

const noise1Bin = 29
const bBin = 31
const noise2Bin = 34
const noise3Bin = 36
const dsBin = 39
const noise4Bin = 42
const noise5Bin = 43
const fsBin = 46
const claimBin = 52
const aBin = 55
const counterClaimBin = 58
const noise6Bin = 60

var inputBuffer *ring[float32]
var inputMutex *sync.Mutex = &sync.Mutex{}

var octaveBuffers [octaves]*ring[float32]

var fftPlan *algofft.PlanRealT[float32, complex64]

var octaveSounds [octaves]OctaveSound

func init() {
	inputBuffer = newRing[float32](1 << 13)

	for i := range octaves {
		octaveBuffers[i] = newRing[float32](1 << 13)
		octaveSounds[i] = OctaveSound{}
	}

	var err error
	fftPlan, err = algofft.NewPlanReal32(fftSize)
	if err != nil {
		log.Fatal(err)
	}
}

func enqueueData(inBuffer []byte, frameCount uint32) {
	inputMutex.Lock()
	for i := range frameCount {
		bits := binary.NativeEndian.Uint16(inBuffer[i*2:])
		amp := float32(int16(bits)) / math.MaxInt16
		inputBuffer.Enqueue(amp)
	}
	inputMutex.Unlock()
}

func analyze() {
	sampleData()

	isDataRemaining := true
	for isDataRemaining {
		for octave := octaves - 1; octave >= 0; octave-- {
			// dequeue data
			timeDomain, success := octaveBuffers[octave].DequeueBatch(fftSize)

			if !success {
				if octave == octaves-1 {
					isDataRemaining = false
				}
				break
			}

			// run fft
			freqDomain, err := applyFFT(timeDomain)

			if err != nil {
				log.Fatal(err)
			}

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
func applyFFT(timeDomain []float32) (freqDomain []complex64, err error) {
	freqDomain = make([]complex64, fftSize/2+1)
	err = fftPlan.Forward(freqDomain, timeDomain)

	return freqDomain, err
}

// downsample data from the input buffer into a buffer ideal for each octave
func sampleData() {
	// take an amount of data that can be downsampled evenly, discard the extra
	usableSampleCount := int(inputBuffer.Size()) & (-1 << (octaves - 1))

	inputMutex.Lock()
	usableSampels, _ := inputBuffer.DequeueBatch(uint64(usableSampleCount))
	inputMutex.Unlock()

	samplings := make([][]float32, octaves)
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
func halfDownSample(input []float32) []float32 {
	var output = make([]float32, len(input)/2)

	for i := range len(input) / 2 {
		output[i] = (input[i*2] + input[i*2+1]) / 2
	}

	return output
}
