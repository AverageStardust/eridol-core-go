package core

import "math/cmplx"

type RawSound struct {
	B            float32
	Ds           float32
	Fs           float32
	A            float32
	Claim        float32
	CounterClaim float32
	Noise        float32
}

const bBin = 28
const dsBin = 35
const fsBin = 42
const claimBin = 47
const aBin = 50
const counterClaimBin = 53

var noiseBins = []int{26, 32, 39, 45, 48, 52, 55}

func catagorizeSound(freqDomain []complex64) RawSound {
	var backgroundNoise float32
	for _, bin := range noiseBins {
		backgroundNoise += complex64Abs(freqDomain[bin])
	}
	backgroundNoise /= float32(len(noiseBins))

	return RawSound{
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

func (a RawSound) Add(b RawSound) RawSound {
	return RawSound{
		B:            a.B + b.B,
		Ds:           a.Ds + b.Ds,
		Fs:           a.Fs + b.Fs,
		A:            a.A + b.A,
		Claim:        a.Claim + b.Claim,
		CounterClaim: a.CounterClaim + b.CounterClaim,
		Noise:        a.Noise + b.Noise,
	}
}

func (a RawSound) Scale(b float32) RawSound {
	return RawSound{
		B:            a.B * b,
		Ds:           a.Ds * b,
		Fs:           a.Fs * b,
		A:            a.A * b,
		Claim:        a.Claim * b,
		CounterClaim: a.CounterClaim * b,
		Noise:        a.Noise * b,
	}
}
