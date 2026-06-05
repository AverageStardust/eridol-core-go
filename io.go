package main

import (
	"encoding/binary"
	"log"
	"time"

	"github.com/gen2brain/malgo"
)

type SoundCallback func(octave int, sound OctaveSound, sampleTime time.Duration)

const sampleRate = 48000

var DoLogging = false

var context *malgo.AllocatedContext
var device *malgo.Device
var userCallback SoundCallback

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

func OnSound(callback SoundCallback) {
	userCallback = callback
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

func sendUserCallback(octave int, sound OctaveSound, sampleTimeSamples uint64) {
	callback := userCallback
	if callback == nil {
		return
	}

	// calculate time to the nearest microsecond (discarding some accuracy to prevent overflows)
	sampleTimeDuration := time.Duration(sampleTimeSamples) * (time.Second / time.Microsecond) / time.Duration(sampleRate) * time.Microsecond
	callback(octave, sound, sampleTimeDuration)
}
