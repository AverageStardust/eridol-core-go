package core

import (
	"log"
	"time"
)

// The callback type implemented by the user if they want processed notes.
type UserCallback func(heard Heard) (stop bool)

// The callback type implemented by the user if they want raw sound data.
type UserCallbackRaw func(heard HeardRaw) (stop bool)

// Runs the callback every time new note data is avalable for at least one octave.
// The first callback argument is the result of sound analysis, as an AnalysisResults[Notes].
// The second callback argument is group of synthesizers to play notes in each octave, as a Choir.
// If you want raw sound data use RunWithRawSound().
// The callback is called off of the main thread/goroutine.
func Run(callback UserCallback) error {
	var rawCallback = newNoteAnalyzer(callback)
	return RunWithRawSound(rawCallback)
}

// Runs the callback every time new sound data is avalable fo r at least one octave.
// Sound data is just raw numbers for the loudness of each note, as well as background noise.
// The first callback argument is the result of basic sound analysis, as an AnalysisResults[RawSound].
// The second callback argument is group of synthesizers to play notes in each octave, as a Choir.
// If you want analyzed note data use Run().
// The callback is called off of the main thread/goroutine.
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
		log.Println("eridolcore: Ready to run user callback.")
	}

	// run untill callback tells us to stop
	for !stopping {
		stopping = callback(<-heardFFT)
	}

	globalChoir.silence()

	if doLogging {
		log.Println("eridolcore: Waiting for synths to go silent.")
	}

	// wait for choir to silence
	time.Sleep(time.Second * (1 / synthFadeSpeed) / sampleRate)

	err = stopIO()

	if doLogging {
		log.Println("eridolcore: Waiting for FFT to stop.")
	}

	// stop fft first, even if we get an error

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

	if doLogging {
		log.Println("eridolcore: Exiting run loop.")
	}

	return err
}
