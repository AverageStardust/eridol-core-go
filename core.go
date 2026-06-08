/*
Module eridol-core-go provides a basic interface to build an eridol implementation.

The following starter program keeps ownership of octave 2:

	core.Run(func(heard core.Heard) (stop bool) {
			// check if someone is playing the claim tone
			isSomeoneClaiming := heard.Octaves[2].Claim

			// if someone is claiming play the counter claim tone
			core.Synth(2).SetNoteImmediately(core.CounterClaim, isSomeoneClaiming)

			// never stop
			return false
		})
*/
package core

import "time"

const OctaveCount = 6

var doLogging = false

func EnableLogging() {
	doLogging = true
}

func DisableLogging() {
	doLogging = false
}

func samplesToDuration(samples uint64) time.Duration {
	return time.Duration(samples) * (time.Second / time.Microsecond) / time.Duration(sampleRate) * time.Microsecond
}

func durationToSamples(duration time.Duration) uint64 {
	return uint64(duration / time.Microsecond * time.Duration(sampleRate) / (time.Second / time.Microsecond))
}
