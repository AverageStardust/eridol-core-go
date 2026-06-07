package core

import (
	"log"

	"github.com/gen2brain/malgo"
)

var context *malgo.AllocatedContext
var device *malgo.Device

// Sets up the sound library to listen and play.
func Init() {
	if context != nil {
		return
	}

	var err error
	context, err = malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		if DoLogging {
			log.Println("miniaudio: ", message)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	if DoLogging {
		log.Println("eridol-core: Created malgo context")
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Duplex)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 1
	deviceConfig.NoPreSilencedOutputBuffer = 1
	deviceConfig.NoClip = 1

	device, err = malgo.InitDevice(context.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: deviceDataCallback,
	})
	if err != nil {
		log.Fatal(err)
	}

	if DoLogging {
		log.Println("eridol-core: Created malgo device")
	}

	err = device.Start()
	if err != nil {
		log.Fatal(err)
	}

	if DoLogging {
		log.Println("eridol-core: Started malgo device")
	}
}

// Closes all oscillators and waits for them to go silent.
// Then closes the sound library.
func Uninit() {
	if isClosingOscillators() {
		return
	}

	closeOscillators(finishUninit)
}

func finishUninit() {
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
			log.Println("eridol-core: Released malgo context")
		}

		context = nil
	}
}

func deviceDataCallback(outBuffer, inBuffer []byte, frameCount uint32) {
	// write output
	writeOscillators(outBuffer, int(frameCount))

	// read input
	enqueueData(inBuffer, frameCount)

	// start analyze
	go analyzeData()
}
