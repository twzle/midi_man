package main

import (
	"context"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	"go.uber.org/zap"
	"log"
	"midi_manipulator/pkg/backlight"
	midiC "midi_manipulator/pkg/config"
	midiD "midi_manipulator/pkg/midi"
	"os"
	"os/signal"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Invalid args, usage: [device-name]")
	}

	stopApp := make(chan os.Signal, 1)
	signal.Notify(stopApp, os.Interrupt)

	logger := zap.Must(zap.NewDevelopment())
	defer midi.CloseDriver()

	signals := make(chan core.Signal)
	deviceConfig := midiC.DeviceConfig{
		DeviceName:        os.Args[1],
		StartupDelay:      200,
		ReconnectInterval: 200,
		Active:            true,
		HoldDelta:         1000,
		Namespace:         "default",
		Controls:          nil,
	}

	backlightConfig, err := backlight.InitConfig("configs/backlight_config.yaml")
	if err != nil {
		logger.Fatal("Unable to init backlight")
	}

	aDevice := midiD.NewDevice(deviceConfig, signals, logger, core.NewCheckManager())
	aDevice.RunDevice(backlightConfig)
	defer aDevice.Stop()

	listenSig, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-listenSig.Done():
				return
			case newSignal := <-signals:
				logger.Info("Received", zap.Any("signal", newSignal))
			}
		}
	}()

	<-stopApp
	cancel()
}
