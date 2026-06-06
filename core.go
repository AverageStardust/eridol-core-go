package main

import (
	"encoding/binary"
	"log"
	"time"

	"github.com/gen2brain/malgo"
)

const OctaveCount = 6
const sampleRate = 48000

var DoLogging = false

var context *malgo.AllocatedContext
var device *malgo.Device
var userSoundCallback SoundCallback

func Init() {
	if context != nil {
		return
	}

	var err error
	context, err = malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		if DoLogging {
			println("miniaudio: ", message)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	if DoLogging {
		println("eridol-core: Created malgo context")
	}

	if device != nil {
		device.Uninit()
		println("eridol-core: Released malgo device")
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Duplex)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 1
	deviceConfig.SampleRate = sampleRate
	deviceConfig.NoPreSilencedOutputBuffer = 1
	deviceConfig.NoClip = 1
	deviceConfig.Alsa.NoMMap = 1

	captureCallbacks := malgo.DeviceCallbacks{
		Data: audioDataCallback,
	}

	device, err := malgo.InitDevice(context.Context, deviceConfig, captureCallbacks)
	if err != nil {
		log.Fatal(err)
	}

	if DoLogging {
		println("eridol-core: Created malgo device")
	}

	err = device.Start()
	if err != nil {
		log.Fatal(err)
	}

	if DoLogging {
		println("eridol-core: Started malgo device")
	}
}

// Runs the callback every time new sound data is avalable for at least one octave.
// Sound data just has raw numbers for the loudness of each note/tone as well as background noise.
// If you want this analyzed into notes use OnNotes().
// The callback is called off of the main thread/goroutine.
func OnSound(callback SoundCallback) {
	userSoundCallback = callback
}

// Runs the callback every time new note data is avalable for at least one octave.
// Note data has already be analized into boolean on/off values for each note.
// If you want raw sound data use OnSound().
// The callback is called off of the main thread/goroutine.
func OnNotes(callback NotesCallback) {
	OnNotesWithThreshhold(callback, 5)
}

// Runs the callback every time new note data is avalable for at least one octave.
// Higher values of signalThreshold be less sensative to quite notes, but less likey to cause false positive detections.
// Note data has already be analized into boolean on/off values for each note.
// If you want raw sound data use OnSound().
// The callback is called off of the main thread/goroutine.
func OnNotesWithThreshhold(callback NotesCallback, signalThreshold float32) {
	userSoundCallback = createNoteAnalyzer(callback, signalThreshold)
}

func Uninit() {
	if device != nil {
		device.Uninit()

		if DoLogging {
			println("eridol-core: Released malgo device")
		}

		device = nil
	}

	if context != nil {
		err := context.Uninit()
		if err != nil {
			log.Fatal(err)
		}

		context.Free()

		if DoLogging {
			println("eridol-core: Released malgo context")
		}

		context = nil
	}
}

func audioDataCallback(outBuffer, inBuffer []byte, frameCount uint32) {
	// read input
	enqueueData(inBuffer, frameCount)
	if DoLogging {
		println("eridol-core: Read ", frameCount, " frames of audio input")
	}

	go analyze()

	if true {
		// fill zeros
		for i := range frameCount {
			bits := uint16(int16(0))
			binary.NativeEndian.PutUint16(outBuffer[i*2:], bits)
		}

		if DoLogging {
			println("eridol-core: Zeroed ", frameCount, " frames of audio output")
		}

	} else {
		// write output
		// for i := range frameCount {
		// 	amp, _ := ouputRing.Dequeue()
		// 	bits := uint16(int16(amp * math.MaxInt16))
		// 	binary.NativeEndian.PutUint16(outBuffer[i*2:], bits)
		// }

		// if DoLogging {
		// 	println("eridol-core: Wrote ", frameCount, " frames of audio output")
		// }
	}
}

func sendUserCallback(octaves [OctaveCount]Sound, analysisTimeSamples uint64) {
	callback := userSoundCallback
	if callback == nil {
		return
	}

	// calculate time to the nearest microsecond (discarding some accuracy to prevent overflows)
	analysisTimeDuration := time.Duration(analysisTimeSamples) * (time.Second / time.Microsecond) / time.Duration(sampleRate) * time.Microsecond
	callback(octaves, analysisTimeDuration)
}
