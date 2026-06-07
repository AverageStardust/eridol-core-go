package core

type Sound struct {
	B            float32
	Ds           float32
	Fs           float32
	A            float32
	Claim        float32
	CounterClaim float32
	Noise        float32
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
