package main

import (
	"fmt"
	"git.miem.hse.ru/hubman/hubman-lib"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"git.miem.hse.ru/hubman/hubman-lib/executor"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	"log"
	"midi_manipulator/pkg/config"
	midiHermophrodite "midi_manipulator/pkg/midi"
	"midi_manipulator/pkg/model"
)

func main() {
	defer midi.CloseDriver()

	systemConfig := &core.SystemConfig{}
	userConfig := &config.UserConfig{}

	err := core.ReadConfig(systemConfig, userConfig)
	if err != nil {
		log.Fatal(fmt.Errorf("error while reading config: %w", err))
	}

	err = userConfig.Validate()
	if err != nil {
		log.Fatal(err)
	}

	err = systemConfig.Validate()
	if err != nil {
		log.Fatal(err)
	}

	setupApp(systemConfig, userConfig)
}

func setupApp(systemConfig *core.SystemConfig, userConfig *config.UserConfig) {
	deviceManager := midiHermophrodite.NewDeviceManager()
	defer deviceManager.Close()

	agentConf := core.AgentConfiguration{
		System:          systemConfig,
		User:            userConfig,
		ParseUserConfig: func(data []byte) (core.Configuration, error) { return config.ParseConfigFromBytes(data) },
	}

	deviceManager.UpdateDevices(userConfig.MidiDevices)
	signals := deviceManager.GetSignals()

	app := hubman.NewAgentApp(
		agentConf,
		hubman.WithManipulator(
			hubman.WithSignal[model.NotePushed](),
			hubman.WithSignal[model.NoteHold](),
			hubman.WithSignal[model.NoteReleased](),
			hubman.WithSignal[model.ControlPushed](),
			hubman.WithChannel(signals),
		),
		hubman.WithExecutor(
			hubman.WithCommand(model.TurnLightOnCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd model.TurnLightOnCommand
					parser(&cmd)
					err := deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					if err != nil {
						log.Println(err)
					}
				}),
			hubman.WithCommand(model.TurnLightOffCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd model.TurnLightOffCommand
					parser(&cmd)
					err := deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					if err != nil {
						log.Println(err)
					}
				}),
			hubman.WithCommand(model.SingleBlinkCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd model.SingleBlinkCommand
					parser(&cmd)
					err := deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					if err != nil {
						log.Println(err)
					}
				}),
			hubman.WithCommand(model.SingleReversedBlinkCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd model.SingleReversedBlinkCommand
					parser(&cmd)
					err := deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					if err != nil {
						log.Println(err)
					}
				}),
			hubman.WithCommand(model.ContinuousBlinkCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd model.ContinuousBlinkCommand
					parser(&cmd)
					err := deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					if err != nil {
						log.Println(err)
					}
				})),
		hubman.WithOnConfigRefresh(func(configuration core.AgentConfiguration) {
			update, _ := configuration.User.([]config.DeviceConfig)
			deviceManager.UpdateDevices(update)
		}),
	)
	<-app.WaitShutdown()
}
