package core

import (
	"encoding/binary"
	"math"
	"time"

	algofft "github.com/cwbudde/algo-fft"

	"github.com/averagestardust/eridol-core-go/internal"
)

// The results of eridol analyzing sound from the microphone.
// Contains boolean values for if each note is declared as playing or not playing
type Heard struct {
	Octaves         [OctaveCount]Notes
	IsOctaveChanged [OctaveCount]bool
	TimeRunning     time.Duration
}

// The results of eridol analyzing sound from the microphone.
// Contains raw numbers for the loudness of each note, as well as background noise.
type HeardRaw struct {
	Octaves         [OctaveCount]RawSound
	IsOctaveChanged [OctaveCount]bool
	TimeRunning     time.Duration
}

const fftSize = 340   // numbers of samples used in an fft run
const fftStride = 170 // step size of sampels between fft runs

func newFFTAnalyzer(heard chan HeardRaw) (run func(soundQuanta), stop chan struct{}, err error) {
	inputRing := internal.NewRing[float32](4096)

	octaveRings := [OctaveCount]*internal.Ring[float32]{}
	for i := range OctaveCount {
		octaveRings[i] = internal.NewRing[float32](4096)
	}

	var fftPlan *algofft.PlanRealT[float32, complex64]
	fftPlan, err = algofft.NewPlanReal32(fftSize)
	if err != nil {
		return nil, nil, err
	}

	// allocated once to stop reallocation
	frequencyDomain := make([]complex64, fftSize/2+1)

	octaveSounds := &[OctaveCount]RawSound{}

	doneEnqueingData := make(chan struct{})
	stop = make(chan struct{})

	go func() {
	fftLoop:
		for {
			select {
			case <-doneEnqueingData:
				sampleFFTData(octaveRings, inputRing)
				analyzeFFTData(heard, octaveSounds, frequencyDomain, fftPlan, octaveRings)
			case <-stop:
				break fftLoop
			}
		}
	}()

	return func(quanta soundQuanta) {
		enqueueFFTData(inputRing, quanta)
		doneEnqueingData <- struct{}{}
	}, stop, nil
}

func enqueueFFTData(inputRing *internal.Ring[float32], quanta soundQuanta) {
	for i := range quanta.frameCount {
		bits := binary.NativeEndian.Uint16(quanta.buffer[i*2:])
		amp := float32(int16(bits)) / math.MaxInt16
		inputRing.Enqueue(amp)
	}
}

// downsample data from the input buffer into a buffer ideal for each octave
func sampleFFTData(octaveRings [OctaveCount]*internal.Ring[float32], inputRing *internal.Ring[float32]) {
	// take an amount of data that can be downsampled evenly, discard the extra
	usableSampleCount := int(inputRing.Size()) & (-1 << (OctaveCount - 1))

	if usableSampleCount == 0 {
		return
	}

	usableSampels, _ := inputRing.DequeueBatch(uint64(usableSampleCount))

	samplings := make([][]float32, OctaveCount)
	samplings[OctaveCount-1] = usableSampels

	// create different down-samplings at an idea sampling rate for each octave
	for i := OctaveCount - 1; i > 0; i-- {
		samplings[i-1] = halfDownSample(samplings[i])
	}

	for i := range OctaveCount {
		octaveRings[i].EnqueueBatch(samplings[i])
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

func analyzeFFTData(heard chan HeardRaw, octaveSounds *[OctaveCount]RawSound, frequencyDomain []complex64, fftPlan *algofft.PlanRealT[float32, complex64], octaveRings [OctaveCount]*internal.Ring[float32]) {
	var analysisSampleTime uint64
	var isOctaveChanged [OctaveCount]bool

outerLoop:
	for {
		for i := range OctaveCount {
			isOctaveChanged[i] = false
		}

		for octave := OctaveCount - 1; octave >= 0; octave-- {
			octaveRing := octaveRings[octave]

			// find the time in samlpes at the midpoint of the current octave's data
			octaveSampleTime := (octaveRing.Tail() + fftSize/2) << (OctaveCount - 1 - octave)

			if octave == OctaveCount-1 {
				// allow the octave with the must frequent data to dictate the time
				analysisSampleTime = octaveSampleTime
			} else if octaveSampleTime > analysisSampleTime {
				// don't do anything that is before it's time
				continue
			}

			// get octave data
			timeDomain, success := octaveRing.PeekBatch(octaveRing.Tail(), fftSize)
			if !success {
				if octave == OctaveCount-1 {
					// out of data for more frequent octave, and thus out of data for all octaves
					break outerLoop
				}
				break
			}

			// move forward in data
			octaveRing.Drop(fftStride)

			// run fft
			fftPlan.Forward(frequencyDomain, timeDomain)

			if visualizedOctave == octave {
				println("------ Octave", octave, "FFT Visualized ------")
				internal.VisualizeFrequencyDomain(frequencyDomain[bBin-2:counterClaimBin+4], 10)
			}

			// categorize bins
			octaveSounds[octave] = catagorizeSound(frequencyDomain)
			isOctaveChanged[octave] = true
		}

		heard <- HeardRaw{
			Octaves:         *octaveSounds,
			IsOctaveChanged: isOctaveChanged,
			TimeRunning:     samplesToDuration(analysisSampleTime),
		}
	}
}
