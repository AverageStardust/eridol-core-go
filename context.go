package main

import (
	"encoding/binary"
	"log"
	"math"

	"github.com/gen2brain/malgo"
)

const sampleRate = 48000

var DoLogging = false

var context *malgo.AllocatedContext
var device *malgo.Device
var inputRing *ring[float64] = newRing[float64](1 << 14)
var ouputRing *ring[float64] = newRing[float64](1 << 14)

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
	for i := range frameCount {
		bits := binary.NativeEndian.Uint16(inBuffer[i*2:])
		amp := float64(int16(bits)) / math.MaxInt16
		inputRing.Enqueue(amp)
	}

	if DoLogging {
		println("eridol-core: Read ", frameCount, " frames of audio input")
	}

	zeroOutput := uint32(ouputRing.Size()) >= frameCount && false

	if zeroOutput {
		// fill zeros
		for i := range frameCount {
			bits := uint16(int16(0))
			binary.NativeEndian.PutUint16(inBuffer[i*2:], bits)
		}

		if DoLogging {
			println("eridol-core: Zeroed ", frameCount, " frames of audio output")
		}

	} else {
		// write output
		for i := range frameCount {
			amp, _ := ouputRing.Dequeue()
			bits := uint16(int16(amp * math.MaxInt16))
			binary.NativeEndian.PutUint16(inBuffer[i*2:], bits)
		}

		if DoLogging {
			println("eridol-core: Wrote ", frameCount, " frames of audio output")
		}
	}
}
