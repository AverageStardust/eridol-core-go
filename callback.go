package core

import (
	"time"
)

type SoundCallback func(octaves [OctaveCount]Sound, runtime time.Duration)

type NotesCallback func(octaves [OctaveCount]Notes, runtime time.Duration)

var userCallback SoundCallback

// Runs the callback every time new sound data is avalable for at least one octave.
// Sound data just has raw numbers for the loudness of each note/tone as well as background noise.
// If you want this analyzed into notes use OnNotes().
// The callback is called off of the main thread/goroutine.
func OnSound(callback SoundCallback) {
	userCallback = callback
}

// Runs the callback every time new note data is avalable for at least one octave.
// Note data has already be analized into boolean on/off values for each note.
// If you want raw sound data use OnSound().
// The callback is called off of the main thread/goroutine.
func OnNotes(callback NotesCallback) {
	OnNotesWithThreshhold(callback, 8)
}

// Runs the callback every time new note data is avalable for at least one octave.
// Higher values of signalThreshold be less sensative to quite notes, but less likey to cause false positive detections.
// Note data has already be analized into boolean on/off values for each note.
// If you want raw sound data use OnSound().
// The callback is called off of the main thread/goroutine.
func OnNotesWithThreshhold(callback NotesCallback, signalThreshold float32) {
	userCallback = createNoteAnalyzer(callback, signalThreshold)
}

func sendUserCallback(octaves [OctaveCount]Sound, analysisSamples uint64) {
	callback := userCallback
	if callback == nil {
		return
	}

	// calculate time to the nearest microsecond (discarding some accuracy to prevent overflows)
	callback(octaves, analysisSamplesToRuntime(analysisSamples))
}

func analysisSamplesToRuntime(analysisSamples uint64) time.Duration {
	return time.Duration(analysisSamples) * (time.Second / time.Microsecond) / time.Duration(sampleRate) * time.Microsecond
}

func runtimeToAnalysisSamples(runtime time.Duration) uint64 {
	return uint64(runtime / time.Microsecond * time.Duration(sampleRate) / (time.Second / time.Microsecond))
}
