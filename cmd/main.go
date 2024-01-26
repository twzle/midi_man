package main

import (
	"fmt"
	"git.miem.hse.ru/hubman/hubman-lib"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"git.miem.hse.ru/hubman/hubman-lib/executor"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	"log"
	"midi_manipulator/pkg/backlight"
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

	setupApp(systemConfig, userConfig)
}

func setupApp(systemConfig *core.SystemConfig, userConfig *config.UserConfig) {
	agentConf := core.AgentConfiguration{
		System:          systemConfig,
		User:            userConfig,
		ParseUserConfig: func(data []byte) (core.Configuration, error) { return config.ParseConfigFromBytes(data) },
	}

	app := core.NewContainer(agentConf.System.Logging)
	deviceManager := midiHermophrodite.NewDeviceManager(app.Logger())
	defer deviceManager.Close()

	backlightConfig, err := backlight.InitConfig("configs/backlight_config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	deviceManager.SetBacklightConfig(backlightConfig)
	signals := deviceManager.GetSignals()

	app.RegisterPlugin(
		hubman.NewAgentPlugin(
			app.Logger(),
			agentConf,
			hubman.WithManipulator(
				hubman.WithSignal[model.NotePushed](),
				hubman.WithSignal[model.NoteHold](),
				hubman.WithSignal[model.NoteReleased](),
				hubman.WithSignal[model.NoteReleasedAfterHold](),
				hubman.WithSignal[model.ControlPushed](),
				hubman.WithSignal[model.NamespaceChanged](),
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
					}),
				hubman.WithCommand(model.SetActiveNamespaceCommand{},
					func(s core.SerializedCommand, p executor.CommandParser) {
						var cmd model.SetActiveNamespaceCommand
						p(&cmd)
						deviceManager.SetActiveNamespace(cmd.Namespace, cmd.Device)
					}),
			),
			hubman.WithOnConfigRefresh(func(configuration core.AgentConfiguration) {
				update, _ := configuration.User.(*config.UserConfig)
				deviceManager.UpdateDevices(update.MidiDevices)
			}),
		),
	)

	deviceManager.UpdateDevices(userConfig.MidiDevices)
	go midiHermophrodite.CheckDevicesHealth(deviceManager)

	<-app.WaitShutdown()
}
