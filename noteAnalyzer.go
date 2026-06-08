package core

func newNoteAnalyzer(callback UserCallback, signalThreshold float32) UserCallbackRaw {
	var consistentNotes [OctaveCount]Notes
	var lastAnalyzedNotes [OctaveCount]Notes
	var lastRawOctaves [OctaveCount]RawSound

	return func(results HeardRaw) bool {
		var analyzedNotes [OctaveCount]Notes

		for octave := range OctaveCount {
			if results.IsOctaveChanged[octave] {
				continue
			}

			// mix last two samples be less sensative to random changes
			averagedOctave := results.Octaves[octave].Scale(0.7).Add(lastRawOctaves[octave].Scale(0.3))
			lastRawOctaves[octave] = results.Octaves[octave]

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

		noteResults := Heard{
			Octaves:         consistentNotes,
			IsOctaveChanged: results.IsOctaveChanged,
			TimeRunning:     results.TimeRunning,
		}

		return callback(noteResults)
	}
}
