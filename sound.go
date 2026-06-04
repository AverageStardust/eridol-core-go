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
