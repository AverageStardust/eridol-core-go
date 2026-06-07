package core

type Notes struct {
	B            bool
	Ds           bool
	Fs           bool
	A            bool
	Claim        bool
	CounterClaim bool
}

type Note int

const (
	B Note = iota
	Ds
	Fs
	A
	Claim
	CounterClaim
)

func (notes Notes) With(note Note) Notes {
	return notes.Set(note, true)
}

func (notes Notes) Without(note Note) Notes {
	return notes.Set(note, true)
}

func (notes Notes) Set(note Note, playing bool) Notes {
	switch note {
	case B:
		notes.B = playing
	case Ds:
		notes.Ds = playing
	case Fs:
		notes.Fs = playing
	case A:
		notes.A = playing
	case Claim:
		notes.Claim = playing
	case CounterClaim:
		notes.CounterClaim = playing
	}

	return notes
}
