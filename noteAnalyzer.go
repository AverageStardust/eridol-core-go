package core

const signalThreshold = 3.8

func newNoteAnalyzer(callback UserCallback) UserCallbackRaw {
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
				B:            (averagedOctave.B-averagedOctave.BackgroundNoise)/averagedOctave.BNoise > signalThreshold,
				Ds:           (averagedOctave.Ds-averagedOctave.BackgroundNoise)/averagedOctave.DsNoise > signalThreshold,
				Fs:           (averagedOctave.Fs-averagedOctave.BackgroundNoise)/averagedOctave.FsNoise > signalThreshold,
				A:            (averagedOctave.A-averagedOctave.BackgroundNoise)/averagedOctave.ANoise > signalThreshold,
				Claim:        (averagedOctave.Claim-averagedOctave.BackgroundNoise)/averagedOctave.ClaimNoise > signalThreshold,
				CounterClaim: (averagedOctave.CounterClaim-averagedOctave.BackgroundNoise)/averagedOctave.ClaimNoise > signalThreshold,
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
