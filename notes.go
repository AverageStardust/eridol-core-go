package core

import (
	"time"
)

type NotesCallback func(octaves [OctaveCount]Notes, analysisTime time.Duration)

type Notes struct {
	B            bool
	Ds           bool
	Fs           bool
	A            bool
	Claim        bool
	CounterClaim bool
}
