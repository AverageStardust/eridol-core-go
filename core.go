package core

import "time"

const OctaveCount = 6

var doLogging = false
var visualizedOctave = -1

// Enables logging from miniaudio and eridol-core.
func EnableLogging() {
	doLogging = true
}

// Disables logging from miniaudio and eridol-core.
func DisableLogging() {
	doLogging = false
}

// Shows the visualization of FFT data for an octave.
// Only one octave can be visable at a time.
func EnableVisualization(octave int) {
	visualizedOctave = octave
}

// Disables the visualization of FFT data.
func DisableVisualization() {
	visualizedOctave = -1
}

func samplesToDuration(samples uint64) time.Duration {
	return time.Duration(samples) * (time.Second / time.Microsecond) / time.Duration(sampleRate) * time.Microsecond
}

func durationToSamples(duration time.Duration) uint64 {
	return uint64(duration / time.Microsecond * time.Duration(sampleRate) / (time.Second / time.Microsecond))
}
