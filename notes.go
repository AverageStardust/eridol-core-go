package core

// A set of notes, including the claim and counter-claim tones.
// Each is set as on (true) or off (false).
type Notes struct {
	B            bool
	Ds           bool
	Fs           bool
	A            bool
	Claim        bool
	CounterClaim bool
}

// An enum for each single note, including the claim and counter-claim tones.
type Note int

const (
	B Note = iota
	Ds
	Fs
	A
	Claim
	CounterClaim
)

// Returns the union of both sets of notes.
// An output note will be playing if it is playing in either input set.
func (a Notes) Union(b Notes) Notes {
	return Notes{
		B:            a.B || b.B,
		Ds:           a.Ds || b.Ds,
		Fs:           a.Fs || b.Fs,
		A:            a.A || b.A,
		Claim:        a.Claim || b.Claim,
		CounterClaim: a.CounterClaim || b.CounterClaim,
	}
}

// Returns the set of notes with one note set to playing.
func (notes Notes) With(note Note) Notes {
	return notes.Set(note, true)
}

// Returns the set of notes with one note set to not playing.
func (notes Notes) Without(note Note) Notes {
	return notes.Set(note, true)
}

// Returns the set with one note set to the playing argument.
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
