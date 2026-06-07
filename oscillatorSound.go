package core

import "math"

type oscillatorSound struct {
	B            float32
	Ds           float32
	Fs           float32
	A            float32
	Claim        float32
	CounterClaim float32
}

func (sound oscillatorSound) sample(time float64, octave int) (amplitude float32) {
	amplitude = sampleSinOscillator(time, 123.47, octave) * sound.B
	amplitude += sampleSinOscillator(time, 155.56, octave) * sound.Ds
	amplitude += sampleSinOscillator(time, 185.00, octave) * sound.Fs
	amplitude += sampleSinOscillator(time, 207.65, octave) * sound.Claim
	amplitude += sampleSinOscillator(time, 220.00, octave) * sound.A
	amplitude += sampleSinOscillator(time, 233.08, octave) * sound.CounterClaim
	return
}

func sampleSinOscillator(time float64, baseFrequency float64, octave int) float32 {
	octaveMultiplier := 1 << octave
	freq := baseFrequency * float64(octaveMultiplier)
	phase := time * freq
	return float32(math.Sin(phase * math.Pi * 2))
}

func (sound *oscillatorSound) stepTo(notes Notes, stepSize float32) {
	stepSoundNoteTo(&sound.B, notes.B, stepSize)
	stepSoundNoteTo(&sound.Ds, notes.Ds, stepSize)
	stepSoundNoteTo(&sound.Fs, notes.Fs, stepSize)
	stepSoundNoteTo(&sound.A, notes.A, stepSize)
	stepSoundNoteTo(&sound.Claim, notes.Claim, stepSize)
	stepSoundNoteTo(&sound.CounterClaim, notes.CounterClaim, stepSize)
}

func stepSoundNoteTo(soundNote *float32, note bool, step float32) {
	if note {
		if *soundNote < 1 {
			*soundNote += step
			if *soundNote > 1 {
				*soundNote = 1
			}
		}
	} else {
		if *soundNote > 0 {
			*soundNote -= step
			if *soundNote < 0 {
				*soundNote = 0
			}
		}
	}
}

func (sound oscillatorSound) maxAmplitude() float32 {
	return sound.B + sound.Ds + sound.Fs + sound.A + sound.Claim + sound.CounterClaim
}
