package core

import (
	"slices"

	"github.com/averagestardust/eridol-core-go/internal"
)

// The raw sound data heard as float point loudnesses.
// Background noise around each note, and on average in the whole octave is also included.
type RawSound struct {
	B                 float32
	Ds                float32
	Fs                float32
	A                 float32
	Claim             float32
	CounterClaim      float32
	BNoise            float32
	DsNoise           float32
	FsNoise           float32
	ANoise            float32
	ClaimNoise        float32
	CounterClaimNoise float32
	BackgroundNoise   float32
}

const bBin = 28
const dsBin = 35
const fsBin = 42
const claimBin = 47
const aBin = 50
const counterClaimBin = 53

var signalContainingBins = []int{28, 35, 36, 42, 47, 49, 50, 52, 53}

func catagorizeSound(freqDomain []complex64) RawSound {
	var backgroundNoise float32 = 0

	for i := bBin - 1; i <= counterClaimBin+2; i++ {
		// skip exact bins to not weight average too high
		if slices.Contains(signalContainingBins, i) {
			continue
		}

		// sum up background noise
		backgroundNoise += internal.C64Abs(freqDomain[i])
	}

	backgroundNoise /= float32(counterClaimBin - bBin + 5 - len(signalContainingBins))

	return RawSound{
		B:            internal.C64Abs(freqDomain[bBin]),
		Ds:           internal.C64Abs(freqDomain[dsBin]),
		Fs:           internal.C64Abs(freqDomain[fsBin]),
		A:            internal.C64Abs(freqDomain[aBin]),
		Claim:        internal.C64Abs(freqDomain[claimBin]),
		CounterClaim: internal.C64Abs(freqDomain[counterClaimBin]),

		BNoise: (internal.C64Abs(freqDomain[bBin-1]) +
			internal.C64Abs(freqDomain[bBin+1])) / 2,
		DsNoise: (internal.C64Abs(freqDomain[dsBin-1]) +
			internal.C64Abs(freqDomain[dsBin+1])) / 2,
		FsNoise: (internal.C64Abs(freqDomain[fsBin-1]) +
			internal.C64Abs(freqDomain[fsBin+1])) / 2,
		ANoise: (internal.C64Abs(freqDomain[aBin-1]) +
			internal.C64Abs(freqDomain[aBin+1])) / 2,
		ClaimNoise: (internal.C64Abs(freqDomain[claimBin-1]) +
			internal.C64Abs(freqDomain[claimBin+1])) / 2,
		CounterClaimNoise: (internal.C64Abs(freqDomain[counterClaimBin-1]) +
			internal.C64Abs(freqDomain[counterClaimBin+1])) / 2,

		BackgroundNoise: backgroundNoise,
	}
}

// Adds two raw noises element-wise and returns the sum.
func (a RawSound) Add(b RawSound) RawSound {
	return RawSound{
		B:            a.B + b.B,
		Ds:           a.Ds + b.Ds,
		Fs:           a.Fs + b.Fs,
		A:            a.A + b.A,
		Claim:        a.Claim + b.Claim,
		CounterClaim: a.CounterClaim + b.CounterClaim,

		BNoise:            a.BNoise + b.BNoise,
		DsNoise:           a.DsNoise + b.DsNoise,
		FsNoise:           a.FsNoise + b.FsNoise,
		ANoise:            a.ANoise + b.ANoise,
		ClaimNoise:        a.ClaimNoise + b.ClaimNoise,
		CounterClaimNoise: a.CounterClaimNoise + b.CounterClaimNoise,

		BackgroundNoise: a.BackgroundNoise + b.BackgroundNoise,
	}
}

// Scales each element in a raw noise and returns it.
func (a RawSound) Scale(b float32) RawSound {
	return RawSound{
		B:            a.B * b,
		Ds:           a.Ds * b,
		Fs:           a.Fs * b,
		A:            a.A * b,
		Claim:        a.Claim * b,
		CounterClaim: a.CounterClaim * b,

		BNoise:            a.BNoise * b,
		DsNoise:           a.DsNoise * b,
		FsNoise:           a.FsNoise * b,
		ANoise:            a.ANoise * b,
		ClaimNoise:        a.ClaimNoise * b,
		CounterClaimNoise: a.CounterClaimNoise * b,

		BackgroundNoise: a.BackgroundNoise * b,
	}
}
