package core

import (
	"time"
)

func createNoteAnalyzer(callback NotesCallback, signalThreshold float32) SoundCallback {
	var consistentNotes [OctaveCount]Notes
	var lastAnalyzedNotes [OctaveCount]Notes
	var lastRawOctaves [OctaveCount]Sound

	return func(rawOctaves [OctaveCount]Sound, analysisTime time.Duration) {
		var analyzedNotes [OctaveCount]Notes

		for octave := range OctaveCount {
			if !IsOctaveUpdated(octave) {
				continue
			}

			// mix last two samples be less sensative to random changes
			averagedOctave := rawOctaves[octave].Scale(0.7).Add(lastRawOctaves[octave].Scale(0.3))
			lastRawOctaves[octave] = rawOctaves[octave]

			// check that
			analyzedNotes[octave] = Notes{
				B:            averagedOctave.B/averagedOctave.Noise > signalThreshold,
				Ds:           averagedOctave.Ds/averagedOctave.Noise > signalThreshold,
				Fs:           averagedOctave.Fs/averagedOctave.Noise > signalThreshold,
				A:            averagedOctave.A/averagedOctave.Noise > signalThreshold,
				Claim:        averagedOctave.Claim/averagedOctave.Noise > signalThreshold,
				CounterClaim: averagedOctave.CounterClaim/averagedOctave.Noise > signalThreshold,
			}

			if analyzedNotes[octave] == lastAnalyzedNotes[octave] || octave == 0 {
				consistentNotes[octave] = analyzedNotes[octave]
			}

			lastAnalyzedNotes[octave] = analyzedNotes[octave]
		}

		callback(consistentNotes, analysisTime)
	}
}
