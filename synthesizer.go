package core

import (
	"sync"
	"time"

	internal "github.com/averagestardust/eridol-core-go/internal"
)

// A handle to planned Synth notes, allows them to be cancled later with CancelPlansFrom().
type SynthHandle struct {
	index uint64
	synth *Synthesizer
}

type synthPlan struct {
	samples uint64
	notes   Notes
	onDone  *[]func()
}

// Creates sound output for one octave by playing the necessary notes.
type Synthesizer struct {
	octave          int
	currentSound    oscillator
	immediateNotes  Notes
	plans           *internal.Ring[synthPlan]
	samplesIntoPlan uint64
	mutex           sync.Mutex
	onDone          []func()
}

const synthFadeSpeed = 0.001

var doneSilencingSynths chan struct{}

// Creates an oscillator to play at one octave.
// Note timing is guaranteed to be consitant with runtime when called inside a callback from OnSound() or OnNotes().
// Timing is undefined outside of callbacks.
// Run .Close() if you are done using the oscillator before the end of the program.
func newSynth(octave int) (synth *Synthesizer) {
	return &Synthesizer{
		octave: octave,
		plans:  internal.NewRing[synthPlan](256),
	}
}

// Silences everything, stopping immediate notes and canceling all planned notes.
func (synth *Synthesizer) Silence() {
	synth.mutex.Lock()
	defer synth.mutex.Unlock()

	synth.immediateNotes = Notes{}
	synth.plans.DropAll()
	synth.samplesIntoPlan = 0
}

// Returns the amount of time remaining in the current series of planned notes.
// Includes any trailing silences that have been planned.
func (synth *Synthesizer) PlanTimeRemaining() time.Duration {
	synth.mutex.Lock()
	defer synth.mutex.Unlock()

	var timeSamples uint64
	for plan := range synth.plans.Iter() {
		timeSamples += plan.samples
	}

	timeSamples -= synth.samplesIntoPlan

	return samplesToDuration(timeSamples)
}

// Returns true if all planned sounds are done playing.
func (synth *Synthesizer) IsAllDone() bool {
	synth.mutex.Lock()
	defer synth.mutex.Unlock()

	return synth.plans.Size() == 0
}

func (synth *Synthesizer) OnAllDone(callback func()) {
	synth.mutex.Lock()
	defer synth.mutex.Unlock()

	if synth.plans.Size() == 0 {
		go callback()
	} else {
		synth.onDone = append(synth.onDone, callback)
	}
}

// Stops a note already set to be playing immediately.
// Does not affect planned notes in any way.
func (synth *Synthesizer) StopNoteImmediately(note Note) {
	synth.SetNoteImmediately(note, false)
}

// Starts playing a single note as soon as possible.
// The note will continue to play regardless of planned notes until stopped.
func (synth *Synthesizer) PlayNoteImmediately(note Note) {
	synth.SetNoteImmediately(note, true)
}

// Sets if an immediate note should be playing or note.
// If a note is set to playing it will continue to play regardless of planned notes until stopped.
func (synth *Synthesizer) SetNoteImmediately(note Note, playing bool) {
	synth.mutex.Lock()
	defer synth.mutex.Unlock()

	synth.immediateNotes = synth.immediateNotes.Set(note, playing)
}

// Sets if every note should immediately be playing or not
// For each note set to playing, it will continue to play regardless of planned notes until stopped.
func (synth *Synthesizer) SetNotesImmediately(notes Notes) {
	synth.mutex.Lock()
	defer synth.mutex.Unlock()

	synth.immediateNotes = notes
}

// Cancels all planned notes and delays.
// Notes started with PlayNoteImmediately() will continue regardless until stopped with StopNoteImmediately().
func (synth *Synthesizer) CancelAllPlans() {
	synth.mutex.Lock()
	defer synth.mutex.Unlock()

	synth.plans.DropAll()
	synth.samplesIntoPlan = 0
}

// Plans a silent delay after any existing planned notes.
// This is the same as running PlanNotes() with an empty set of notes.
// Notes started with PlayNoteImmediately() will continue regardless until stopped with StopNoteImmediately().
// Returns a handle to cancel this later.
func (synth *Synthesizer) PlanDelay(duration time.Duration) (handle SynthHandle) {
	return synth.PlanNotes(Notes{}, duration)
}

// Plans to play notes for a set duration
// Will run after any already planned notes, or otherwise as soon as possible.
// Notes started with PlayNoteImmediately() will continue regardless until stopped with StopNoteImmediately().
// Returns a handle to cancel this later.
func (synth *Synthesizer) PlanNotes(notes Notes, duration time.Duration) (handle SynthHandle) {
	synth.mutex.Lock()
	defer synth.mutex.Unlock()

	handle = SynthHandle{synth.plans.Head(), synth}

	synth.plans.Enqueue(synthPlan{
		samples: durationToSamples(duration),
		notes:   notes,
		onDone:  &[]func(){},
	})

	return handle
}

func (synth *Synthesizer) sample(timeSamples uint64) (amplitude float32, maxAmplitude float32) {
	notes := synth.immediateNotes

	plan, success := synth.advancePlan()
	if success {
		notes = notes.Union(plan.notes)
	}

	synth.currentSound = synth.currentSound.stepTowards(notes, synthFadeSpeed)

	time := float64(timeSamples) / float64(sampleRate)

	amplitude = synth.currentSound.sample(time, synth.octave)
	maxAmplitude = synth.currentSound.maxAmplitude()

	return
}

func (synth *Synthesizer) advancePlan() (plan synthPlan, success bool) {
	// dequeue finished plans until we find the current plan
	for {
		plan, success = synth.plans.Peek(synth.plans.Tail())

		if !success {
			for _, callback := range synth.onDone {
				go callback()
			}
			synth.onDone = []func(){}
			return
		}

		if plan.samples <= synth.samplesIntoPlan {
			for _, callback := range *plan.onDone {
				go callback()
			}
			*plan.onDone = []func(){}

			synth.plans.Dequeue()
			synth.samplesIntoPlan = 0
		} else {
			break
		}
	}

	// advance plan
	synth.samplesIntoPlan++

	return
}

// Cancels any planned notes from the handle onwards.
// Notes started with synth.PlayNoteImmediately() will continue regardless until stopped with synthStopNoteImmediately().
func (handle SynthHandle) CancelPlansFrom() {
	synth := handle.synth

	synth.mutex.Lock()
	defer synth.mutex.Unlock()

	synth.plans.DropFrom(uint64(handle.index))

	if synth.plans.Size() == 0 {
		synth.samplesIntoPlan = 0
	}
}

// Returns true if all the handle's sound is done playing.
func (handle SynthHandle) IsDone() bool {
	synth := handle.synth

	synth.mutex.Lock()
	defer synth.mutex.Unlock()

	return synth.plans.Tail() > handle.index
}

func (handle SynthHandle) OnDone(callback func()) bool {
	synth := handle.synth

	synth.mutex.Lock()
	defer synth.mutex.Unlock()

	if synth.plans.Tail() > handle.index {
		go callback()
		return true
	}

	element, success := synth.plans.Peek(handle.index)

	if success {
		*element.onDone = append(*element.onDone, callback)
	}

	return success
}
