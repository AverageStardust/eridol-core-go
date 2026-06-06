package main

import (
	"encoding/binary"
	"log"
	"math"
	"sync"

	algofft "github.com/cwbudde/algo-fft"
)

const fftSize = 374   // numbers of samples used in an fft run
const fftStride = 187 // step size of sampels between fft runs

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
var octaveBuffers [OctaveCount]*ring[float32]
var fftPlan *algofft.PlanRealT[float32, complex64]
var analysisTime uint64
var octaveSounds [OctaveCount]Sound
var octaveChanged [OctaveCount]bool
var inputMutex *sync.Mutex = &sync.Mutex{}
var analyzeMutex *sync.Mutex = &sync.Mutex{}

// return if the octave's sound/notes changed before the last callback
// higher frequency octaves are updated more often because their notes can be detected faster
func IsOctaveUpdated(octave int) bool {
	return octaveChanged[octave]
}

func init() {
	inputBuffer = newRing[float32](1 << 13)

	for i := range OctaveCount {
		octaveBuffers[i] = newRing[float32](1 << 13)
		octaveSounds[i] = Sound{}
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
	analyzeMutex.Lock()
	sampleData()

outerLoop:
	for {
		for i := range OctaveCount {
			octaveChanged[i] = false
		}

		for octave := OctaveCount - 1; octave >= 0; octave-- {
			buffer := octaveBuffers[octave]
			sampleTime := (buffer.Tail() + fftSize/2) << (OctaveCount - 1 - octave)

			if octave == OctaveCount-1 {
				// allow the octave with the must frequent data to dictate the time
				analysisTime = sampleTime
			} else if sampleTime > analysisTime {
				// don't do anything that is before it's time
				continue
			}

			// get data
			timeDomain, success := buffer.PeekBatch(buffer.Tail(), fftSize)

			if !success {
				if octave == OctaveCount-1 {
					// out of data
					break outerLoop
				}
				break
			}

			// move forward in data
			buffer.Drop(fftStride)

			// run fft
			freqDomain, err := applyFFT(timeDomain)
			if err != nil {
				log.Fatal(err)
			}

			if DoLogging {
				println("eridol-core: Ran fft for octave ", octave)
			}

			// categorize bins
			octaveSounds[octave] = newSound(freqDomain)
			octaveChanged[octave] = true
		}

		sendUserCallback(octaveSounds, analysisTime)
	}

	analyzeMutex.Unlock()
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
	usableSampleCount := int(inputBuffer.Size()) & (-1 << (OctaveCount - 1))

	inputMutex.Lock()
	usableSampels, _ := inputBuffer.DequeueBatch(uint64(usableSampleCount))
	inputMutex.Unlock()

	samplings := make([][]float32, OctaveCount)
	samplings[OctaveCount-1] = usableSampels

	// create different down-samplings at an idea sampling rate for each octave
	for i := OctaveCount - 1; i > 0; i-- {
		samplings[i-1] = halfDownSample(samplings[i])
	}

	for i := range OctaveCount {
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
