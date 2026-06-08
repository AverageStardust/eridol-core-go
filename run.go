package core

import "time"

// The callback type implemented by the user if they want processed notes.
type UserCallback func(heard Heard) (stop bool)

// The callback type implemented by the user if they want raw sound data.
type UserCallbackRaw func(heard HeardRaw) (stop bool)

const noteThreshold = 6

// Runs the callback every time new note data is avalable for at least one octave.
// The first callback argument is the result of sound analysis, as an AnalysisResults[Notes].
// The second callback argument is group of synthesizers to play notes in each octave, as a Choir.
// If you want raw sound data use RunWithRawSound().
// The callback is called off of the main thread/goroutine.
func Run(callback UserCallback) error {
	var rawCallback = newNoteAnalyzer(callback, noteThreshold)
	return RunWithRawSound(rawCallback)
}

// Runs the callback every time new sound data is avalable fo r at least one octave.
// Sound data is just raw numbers for the loudness of each note, as well as background noise.
// The first callback argument is the result of basic sound analysis, as an AnalysisResults[RawSound].
// The second callback argument is group of synthesizers to play notes in each octave, as a Choir.
// If you want analyzed note data use Run().
// The callback is called off of the main thread/goroutine.
func RunWithRawSound(callback UserCallbackRaw) error {
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

	stopIO, err := newIO(wrappedRunFFT)
	if err != nil {
		return err
	}

	// run untill callback tells us to stop
	for !stopping {
		stopping = callback(<-heardFFT)
	}

	globalChoir.silence()

	// wait for choir to silence
	time.Sleep(time.Second * (1 / synthFadeSpeed) / sampleRate)

	err = stopIO()

	// stop fft first, even if we get an error
	stopFFT <- struct{}{}

	return err
}
