package core

import (
	"math/cmplx"
)

const bBin = 21
const dsBin = 26
const fsBin = 31
const claimBin = 35
const aBin = 37
const counterClaimBin = 39

const noise1Bin = 19
const noise2Bin = 23
const noise3Bin = 24
const noise4Bin = 33
const noise5Bin = 41

func catagorizeSound(freqDomain []complex64) Sound {
	backgroundNoise := (complex64Abs(freqDomain[noise1Bin]) +
		complex64Abs(freqDomain[noise2Bin]) +
		complex64Abs(freqDomain[noise3Bin]) +
		complex64Abs(freqDomain[noise4Bin]) +
		complex64Abs(freqDomain[noise5Bin])) / 5.0

	return Sound{
		B:            complex64Abs(freqDomain[bBin]),
		Ds:           complex64Abs(freqDomain[dsBin]),
		Fs:           complex64Abs(freqDomain[fsBin]),
		Claim:        complex64Abs(freqDomain[claimBin]),
		A:            complex64Abs(freqDomain[aBin]),
		CounterClaim: complex64Abs(freqDomain[counterClaimBin]),
		Noise:        backgroundNoise,
	}
}

func complex64Abs(n complex64) float32 {
	// ughhh go please
	return float32(cmplx.Abs(complex128(n)))
}
