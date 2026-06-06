package main

import (
	"math/cmplx"
	"time"
)

type SoundCallback func(octaves [OctaveCount]Sound, analysisTime time.Duration)

type Sound struct {
	B            float32
	Ds           float32
	Fs           float32
	A            float32
	Claim        float32
	CounterClaim float32
	Noise        float32
}

func newSound(freqDomain []complex64) Sound {
	backgroundNoise := (complex64Abs(freqDomain[noise1Bin]) +
		complex64Abs(freqDomain[noise2Bin]) +
		complex64Abs(freqDomain[noise3Bin]) +
		complex64Abs(freqDomain[noise4Bin]) +
		complex64Abs(freqDomain[noise5Bin]) +
		complex64Abs(freqDomain[noise6Bin])) / 6.0 * 0.8

	return Sound{
		B:            complex64Abs(freqDomain[bBin]) * 0.51,
		Ds:           complex64Abs(freqDomain[dsBin]) * 0.78,
		Fs:           complex64Abs(freqDomain[fsBin]) * 0.92,
		Claim:        complex64Abs(freqDomain[claimBin]) * 1.04,
		A:            complex64Abs(freqDomain[aBin]) * 1.1,
		CounterClaim: complex64Abs(freqDomain[counterClaimBin]) * 1.36,
		Noise:        backgroundNoise,
	}
}

func complex64Abs(n complex64) float32 {
	// ughhh go please
	return float32(cmplx.Abs(complex128(n)))
}

func (a Sound) Add(b Sound) Sound {
	return Sound{
		B:            a.B + b.B,
		Ds:           a.Ds + b.Ds,
		Fs:           a.Fs + b.Fs,
		A:            a.A + b.A,
		Claim:        a.Claim + b.Claim,
		CounterClaim: a.CounterClaim + b.CounterClaim,
		Noise:        a.Noise + b.Noise,
	}
}

func (a Sound) Scale(b float32) Sound {
	return Sound{
		B:            a.B * b,
		Ds:           a.Ds * b,
		Fs:           a.Fs * b,
		A:            a.A * b,
		Claim:        a.Claim * b,
		CounterClaim: a.CounterClaim * b,
		Noise:        a.Noise * b,
	}
}
