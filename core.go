/*
eridol-core-go provides a basic interface to build a eridol implementation.
*/
package core

import "time"

const OctaveCount = 6

func samplesToDuration(samples uint64) time.Duration {
	return time.Duration(samples) * (time.Second / time.Microsecond) / time.Duration(sampleRate) * time.Microsecond
}

func durationToSamples(duration time.Duration) uint64 {
	return uint64(duration / time.Microsecond * time.Duration(sampleRate) / (time.Second / time.Microsecond))
}
