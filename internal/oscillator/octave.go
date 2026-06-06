package oscillator

import "math"

type octaveSound struct {
	B            float32
	Ds           float32
	Fs           float32
	A            float32
	Claim        float32
	CounterClaim float32
}

func (sound octaveSound) Sample(time float64, octave int) (amplitude float32) {
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

func (sound octaveSound) TotalAmplitude() float32 {
	return sound.B + sound.Ds + sound.Fs + sound.A + sound.Claim + sound.CounterClaim
}
