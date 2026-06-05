package main

import "math/cmplx"

type OctaveSound struct {
	B            float64
	Ds           float64
	Fs           float64
	A            float64
	Claim        float64
	CounterClaim float64
	Noise        float64
}

func newOctaveSound(freqDomain []complex128) OctaveSound {
	backgroundNoise := (cmplx.Abs(freqDomain[noise1Bin]) +
		cmplx.Abs(freqDomain[noise2Bin]) +
		cmplx.Abs(freqDomain[noise3Bin]) +
		cmplx.Abs(freqDomain[noise4Bin]) +
		cmplx.Abs(freqDomain[noise5Bin]) +
		cmplx.Abs(freqDomain[noise6Bin])) / 6.0 * 0.8

	return OctaveSound{
		B:            cmplx.Abs(freqDomain[bBin]) * 0.51,
		Ds:           cmplx.Abs(freqDomain[dsBin]) * 0.78,
		Fs:           cmplx.Abs(freqDomain[fsBin]) * 0.92,
		Claim:        cmplx.Abs(freqDomain[claimBin]) * 1.04,
		A:            cmplx.Abs(freqDomain[aBin]) * 1.1,
		CounterClaim: cmplx.Abs(freqDomain[counterClaimBin]) * 1.36,
		Noise:        backgroundNoise,
	}
}

func (a OctaveSound) Add(b OctaveSound) OctaveSound {
	return OctaveSound{
		B:            a.B + b.B,
		Ds:           a.Ds + b.Ds,
		Fs:           a.Fs + b.Fs,
		A:            a.A + b.A,
		Claim:        a.Claim + b.Claim,
		CounterClaim: a.CounterClaim + b.CounterClaim,
		Noise:        a.Noise + b.Noise,
	}
}

func (a OctaveSound) Scale(b float64) OctaveSound {
	return OctaveSound{
		B:            a.B * b,
		Ds:           a.Ds * b,
		Fs:           a.Fs * b,
		A:            a.A * b,
		Claim:        a.Claim * b,
		CounterClaim: a.CounterClaim * b,
		Noise:        a.Noise * b,
	}
}
