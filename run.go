package core

import (
	"fmt"
	"log"
	"time"
)

// The callback type implemented by the user if they want processed notes.
type UserCallback func(heard Heard) (stop bool)

// The callback type implemented by the user if they want raw sound data.
type UserCallbackRaw func(heard HeardRaw) (stop bool)

const defaultSignalThreshold = 4

// Runs the callback every time new note data is avalable for at least one octave.
// Will not return until the callback returns true to stop.
// If you want raw sound data use RunWithRawSound().
// If you want to set the threshold for note detection, use RunWithThreshold().
func Run(callback UserCallback) error {
	return RunWithThreshold(callback, 4)
}

// Runs the callback every time new note data is avalable for at least one octave.
// The sensativity to detecting notes can be set with signalThreshold.
// If the sensativity is too low it may detect false positives, if it's too high it won't be able to hear quite notes.
// Will not return until the callback returns true to stop.
// The default threshould used by Run() is 4, use that if you want a safe value.
func RunWithThreshold(callback UserCallback, signalThreshold float32) error {
	var rawCallback = newNoteAnalyzer(callback, signalThreshold)
	return RunWithRawSound(rawCallback)
}

// Runs the callback every time new sound data is avalable fo r at least one octave.
// Sound data is just raw numbers for the loudness of each note, as well as background noise.
// Will not return until the callback returns true to stop.
// If you want analyzed note data use Run().
func RunWithRawSound(callback UserCallbackRaw) error {
	if doLogging {
		log.Println("eridolcore: Entering run loop.")
	}

	heardFFT := make(chan HeardRaw)

	runFFT, stopFFT, err := newFFTAnalyzer(heardFFT)
	if err != nil {
		return err
	}

	// wrap fft so we can stop analysis when shutting down
	stopping := false
	wrappedRunFFT := func(quanta soundQuanta) {
		if !stopping {
			runFFT(quanta)
		}
	}

	if doLogging {
		log.Println("eridolcore: Waiting for IO to start.")
	}
	stopIO, err := newIO(wrappedRunFFT)
	if err != nil {
		return err
	}

	if doLogging {
		log.Println("eridolcore: Waiting to run user callback.")
	}

	// run until callback tells us to stop
	for !stopping {
		if doLogging {
			log.Println("eridolcore: Running user callback.")
		}
		stopping = callback(<-heardFFT)
	}

	silenceChoirAndWait()

	err = stopIO()

	ignoreUntilFFTStops(stopFFT, heardFFT)

	if doLogging {
		log.Println("eridolcore: Exiting run loop.")
	}

	return err
}

func silenceChoirAndWait() {
	globalChoir.silence()

	// wait for choir to silence
	faidTime := time.Second*(1/synthFadeSpeed)/time.Duration(sampleRate) + time.Millisecond*10
	if doLogging {
		fmt.Printf("eridolcore: Waiting for %v synths to go silent.\n", faidTime)
	}
	time.Sleep(time.Second * (1 / synthFadeSpeed) / time.Duration(sampleRate))
}

func ignoreUntilFFTStops(stopFFT chan struct{}, heardFFT chan HeardRaw) {
	if doLogging {
		log.Println("eridolcore: Waiting for FFT to stop.")
	}
	// stop fft first, even if we just got an error
loop:
	for {
		select {
		case stopFFT <- struct{}{}:
			break loop
		case <-heardFFT:
			// discard
		}
	}
}
