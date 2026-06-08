package core

import (
	"math"

	"github.com/averagestardust/eridol-core-go/internal"
)

type oscillator struct {
	B            float32
	Ds           float32
	Fs           float32
	A            float32
	Claim        float32
	CounterClaim float32
}

func (osc oscillator) sample(time float64, octave int) (amplitude float32) {
	octaveMultiplier := float64(uint(1) << octave)

	amplitude = oscillatorWave(time, 123.47*octaveMultiplier) * osc.B
	amplitude += oscillatorWave(time, 155.56*octaveMultiplier) * osc.Ds
	amplitude += oscillatorWave(time, 185.00*octaveMultiplier) * osc.Fs
	amplitude += oscillatorWave(time, 207.65*octaveMultiplier) * osc.Claim
	amplitude += oscillatorWave(time, 220.00*octaveMultiplier) * osc.A
	amplitude += oscillatorWave(time, 233.08*octaveMultiplier) * osc.CounterClaim

	return
}

func oscillatorWave(time float64, frequency float64) float32 {
	phase := time * frequency
	return float32(math.Sin(phase * math.Pi * 2))
}

func (osc oscillator) maxAmplitude() float32 {
	return osc.B + osc.Ds + osc.Fs + osc.A + osc.Claim + osc.CounterClaim
}

func (osc oscillator) stepTowards(notes Notes, step float32) oscillator {
	osc.B = stepSoundToNote(osc.B, notes.B, step)
	osc.Ds = stepSoundToNote(osc.Ds, notes.Ds, step)
	osc.Fs = stepSoundToNote(osc.Fs, notes.Fs, step)
	osc.A = stepSoundToNote(osc.A, notes.A, step)
	osc.Claim = stepSoundToNote(osc.Claim, notes.Claim, step)
	osc.CounterClaim = stepSoundToNote(osc.CounterClaim, notes.CounterClaim, step)

	return osc
}

func stepSoundToNote(sound float32, note bool, step float32) float32 {
	if note {
		return internal.StepTowards(sound, 1, step)
	} else {
		return internal.StepTowards(sound, 0, step)
	}
}
