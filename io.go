package core

import (
	"log"

	"github.com/gen2brain/malgo"
)

type soundQuanta struct {
	buffer     []byte
	frameCount uint32
}

const sampleRate = 48000

func newIO(runFFT func(soundQuanta)) (stop func() error, err error) {
	context, err := initContext(false)
	if err != nil {
		return nil, err
	}

	device, err := initDevice(context.Context, func(outBuffer, inBuffer []byte, frameCount uint32) {
		// start fft analyzer
		runFFT(soundQuanta{
			buffer:     inBuffer,
			frameCount: frameCount,
		})

		// write synth output
		choir.writeTo(soundQuanta{
			buffer:     outBuffer,
			frameCount: frameCount,
		})
	})

	return func() error {
		device.Uninit()
		return uninitContext(context)
	}, nil
}

func initContext(doLogging bool) (context *malgo.AllocatedContext, err error) {
	return malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		if doLogging {
			log.Println("miniaudio: ", message)
		}
	})
}

func initDevice(context malgo.Context, dataCallback malgo.DataProc) (device *malgo.Device, err error) {
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Duplex)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 1
	deviceConfig.SampleRate = sampleRate
	deviceConfig.NoPreSilencedOutputBuffer = 1
	deviceConfig.NoClip = 1

	device, err = malgo.InitDevice(context, deviceConfig, malgo.DeviceCallbacks{Data: dataCallback})
	if err != nil {
		return nil, err
	}

	err = device.Start()
	return
}

func uninitContext(context *malgo.AllocatedContext) (err error) {
	err = context.Uninit()
	if err != nil {
		return err
	}

	context.Free()

	return nil
}
