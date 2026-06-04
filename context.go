package "github.com/averagestardust/eridol-core-go"

import (
	"container/ring"
	"log"

	"github.com/gen2brain/malgo"
)

const sampleRate = 48000

var DoLogging = false

var context *malgo.AllocatedContext
var device *malgo.Device
var osc *oscillator
var input *ring.Ring

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
	deviceConfig.Capture.Format = malgo.FormatF32
	deviceConfig.Capture.Channels = 1
	deviceConfig.Playback.Format = malgo.FormatF32
	deviceConfig.Playback.Channels = 1
	deviceConfig.SampleRate = sampleRate
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

	osc = newOscillator(sampleRate)
}

func Uninit() {
	osc = nil

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
	if DoLogging {
		println("eridol-core: Received ", frameCount, " frames of audio data")
	}

	osc.write(outBuffer, int(frameCount))
}
