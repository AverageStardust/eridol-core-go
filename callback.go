package core

import (
	"time"
)

var userCallback SoundCallback

func sendUserCallback(octaves [OctaveCount]Sound, analysisTimeSamples uint64) {
	callback := userCallback
	if callback == nil {
		return
	}

	// calculate time to the nearest microsecond (discarding some accuracy to prevent overflows)
	analysisTimeDuration := time.Duration(analysisTimeSamples) * (time.Second / time.Microsecond) / time.Duration(sampleRate) * time.Microsecond
	callback(octaves, analysisTimeDuration)
}

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
