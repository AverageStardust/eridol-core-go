package main

import (
	"time"
)

type NotesCallback func(octaves [OctaveCount]Notes, analysisTime time.Duration)

type Notes struct {
	B            bool
	Ds           bool
	Fs           bool
	A            bool
	Claim        bool
	CounterClaim bool
}

func createNoteAnalyzer(callback NotesCallback, signalThreshold float32) SoundCallback {
	var consistentNotes [OctaveCount]Notes
	var lastAnalyzedNotes [OctaveCount]Notes
	var lastRawOctaves [OctaveCount]Sound

	return func(rawOctaves [OctaveCount]Sound, analysisTime time.Duration) {
		var analyzedNotes [OctaveCount]Notes

		for i := range OctaveCount {
			if !IsOctaveUpdated(i) {
				continue
			}

			// average over time to be less sensative to minor changes
			averagedOctave := rawOctaves[i].Add(lastRawOctaves[i]).Scale(0.5)
			lastRawOctaves[i] = rawOctaves[i]

			// check that
			analyzedNotes[i] = Notes{
				B:            averagedOctave.B/averagedOctave.Noise > signalThreshold,
				Ds:           averagedOctave.Ds/averagedOctave.Noise > signalThreshold,
				Fs:           averagedOctave.Fs/averagedOctave.Noise > signalThreshold,
				A:            averagedOctave.A/averagedOctave.Noise > signalThreshold,
				Claim:        averagedOctave.Claim/averagedOctave.Noise > signalThreshold,
				CounterClaim: averagedOctave.CounterClaim/averagedOctave.Noise > signalThreshold,
			}

			if analyzedNotes[i] == lastAnalyzedNotes[i] {
				consistentNotes[i] = analyzedNotes[i]
			}

			lastAnalyzedNotes[i] = analyzedNotes[i]
		}

		callback(consistentNotes, analysisTime)
	}
}
