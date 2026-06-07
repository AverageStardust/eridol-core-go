package core

import (
	"log"
	"sync"

	internal "github.com/averagestardust/eridol-core-go/internal"
)

type oscillatorPlan struct {
	timeSamples uint64
	targetNotes Notes
}

type Oscillator struct {
	octave       int
	currentSound oscillatorSound
	targetNotes  Notes
	plans        *internal.Ring[*oscillatorPlan]
	alive        bool
}

var oscillators = internal.Set[*Oscillator]{}
var oscillatorMutex = sync.Mutex{}

// Creates an oscillator to play at one octave.
// Note timing is guaranteed to be consitant with runtime when called inside a callback from OnSound() or OnNotes().
// Timing is undefined outside of callbacks.
// Run .Close() if you are done using the oscillator before the end of the program.
func NewOscillator(octave int) (osc *Oscillator) {
	osc = &Oscillator{
		octave: octave,
		plans:  internal.NewRing[*oscillatorPlan](4),
		alive:  true,
	}

	if isClosingOscillators() {
		log.Println("eridol-core: Warning! Created oscillator while Uninit() is still closing. Oscillator will be silent.")
		osc.alive = false
		return osc
	}

	oscillatorMutex.Lock()
	oscillators.Add(osc)
	oscillatorMutex.Unlock()

	return osc
}

// Starts one note playing.
func (osc *Oscillator) PlayNote(note Note) {
	osc.SetNote(note, true)
}

// Stops one note from playing.
func (osc *Oscillator) StopNote(note Note) {
	osc.SetNote(note, false)
}

// Sets at single note to playing or not playing.
func (osc *Oscillator) SetNote(note Note, playing bool) {
	oscillatorMutex.Lock()
	defer oscillatorMutex.Unlock()

	plannedTime := fftSamples + oscillatorLatency

	if plannedTime < oscillatorSamples {
		log.Println("eridol-core: Warning! Missed note change because code is too slow.")
		return
	}

	nextPlan, success := osc.plans.Peek(osc.plans.Tail())

	if success {
		if nextPlan.timeSamples < fftSamples {
			// add new plan after
			osc.plans.Enqueue(&oscillatorPlan{
				timeSamples: plannedTime,
				targetNotes: nextPlan.targetNotes.Set(note, playing),
			})
		} else {
			// update plan
			nextPlan.targetNotes = nextPlan.targetNotes.Set(note, playing)
		}
	} else {
		// add first plan
		osc.plans.Enqueue(&oscillatorPlan{
			timeSamples: plannedTime,
			targetNotes: osc.targetNotes.Set(note, playing),
		})
	}
}

// Stops all notes from playing on oscillator.
func (osc *Oscillator) Stop() {
	osc.Play(Notes{})
}

// Plays some set of notes on oscillator.
func (osc *Oscillator) Play(notes Notes) {
	oscillatorMutex.Lock()
	defer oscillatorMutex.Unlock()

	plannedTime := fftSamples + oscillatorLatency

	if plannedTime < oscillatorSamples {
		log.Println("eridol-core: Warning! Missed note change because code is too slow.")
		return
	}

	nextPlan, success := osc.plans.Peek(osc.plans.Tail())

	if !success || nextPlan.timeSamples < fftSamples {
		// add first plan or add new plan after
		osc.plans.Enqueue(&oscillatorPlan{
			timeSamples: plannedTime,
			targetNotes: notes,
		})
	} else {
		// update plan
		nextPlan.targetNotes = notes
	}
}

// Closes the oscillator so it will stop making sound.
func (osc *Oscillator) Close() {
	oscillatorMutex.Lock()
	defer oscillatorMutex.Unlock()

	osc.alive = false
	osc.plans.DropAll()

	osc.plans.Enqueue(&oscillatorPlan{
		timeSamples: oscillatorSamples,
		targetNotes: Notes{},
	})
}

// Checks if the oscillator has been closed.
func (osc *Oscillator) IsClosed() bool {
	return !osc.alive
}

func (osc *Oscillator) sample() (amplitude float32, maxAmplitude float32) {
	nextPlan, success := osc.plans.Peek(osc.plans.Tail())

	if success && oscillatorSamples >= nextPlan.timeSamples {
		osc.targetNotes = nextPlan.targetNotes
		osc.plans.Dequeue()
	}

	time := float64(oscillatorSamples) / float64(sampleRate)

	osc.currentSound.stepTo(osc.targetNotes, 0.001)

	amplitude = osc.currentSound.sample(time, osc.octave)
	maxAmplitude = osc.currentSound.maxAmplitude()

	if !osc.alive && maxAmplitude == 0 {
		oscillators.Delete(osc)
		if len(oscillators) == 0 {
			if afterLastOscillatorGoesSilent != nil {
				afterLastOscillatorGoesSilent()
			}
		}
	}

	return
}
