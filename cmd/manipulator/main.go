package main

import (
	"go.uber.org/zap"
	"log"
	"midi_manipulator/pkg/backlight"
	"midi_manipulator/pkg/config"
	midiHermophrodite "midi_manipulator/pkg/midi"
	"midi_manipulator/pkg/model"

	"git.miem.hse.ru/hubman/hubman-lib"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"git.miem.hse.ru/hubman/hubman-lib/executor"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

// Application entry point
func main() {
	defer midi.CloseDriver()

	systemConfig := &core.SystemConfig{}
	userConfig := &config.UserConfig{}

	err := core.ReadConfig(systemConfig, userConfig)
	if err != nil {
		log.Fatalf("error while reading config: %e", err)
	}

	setupApp(systemConfig, userConfig)
}

// Function with initialization of application core with configuration data. Ends with termination of process.
func setupApp(systemConfig *core.SystemConfig, userConfig *config.UserConfig) {
	agentConf := core.AgentConfiguration{
		System:          systemConfig,
		User:            userConfig,
		ParseUserConfig: func(data []byte) (core.Configuration, error) { return config.ParseConfigFromBytes(data) },
	}

	app := core.NewContainer(agentConf.System.Logging)
	logger := app.Logger()
	checkManager := core.NewCheckManager()

	deviceManager := midiHermophrodite.NewDeviceManager(logger, checkManager)
	defer deviceManager.Close()

	backlightConfig, err := backlight.InitConfig("configs/backlight_config.yaml")
	if err != nil {
		logger.Fatal("can't init backlight config", zap.Error(err))
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
					func(command core.SerializedCommand, parser executor.CommandParser) error {
						var cmd model.TurnLightOnCommand
						parser(&cmd)
						return deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					}),
				hubman.WithCommand(model.TurnLightOffCommand{},
					func(command core.SerializedCommand, parser executor.CommandParser) error {
						var cmd model.TurnLightOffCommand
						parser(&cmd)
						return deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					}),
				hubman.WithCommand(model.SingleBlinkCommand{},
					func(command core.SerializedCommand, parser executor.CommandParser) error {
						var cmd model.SingleBlinkCommand
						parser(&cmd)
						return deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					}),
				hubman.WithCommand(model.SingleReversedBlinkCommand{},
					func(command core.SerializedCommand, parser executor.CommandParser) error {
						var cmd model.SingleReversedBlinkCommand
						parser(&cmd)
						return deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					}),
				hubman.WithCommand(model.ContinuousBlinkCommand{},
					func(command core.SerializedCommand, parser executor.CommandParser) error {
						var cmd model.ContinuousBlinkCommand
						parser(&cmd)
						return deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					}),
				hubman.WithCommand(model.SetActiveNamespaceCommand{},
					func(s core.SerializedCommand, parser executor.CommandParser) error {
						var cmd model.SetActiveNamespaceCommand
						parser(&cmd)
						return deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					}),
				hubman.WithCommand(model.StartBlinkingCommand{}, func(s core.SerializedCommand, parser executor.CommandParser) error {
					var cmd model.StartBlinkingCommand
					parser(&cmd)
					return deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
				}),
				hubman.WithCommand(model.StopBlinkingCommand{}, func(s core.SerializedCommand, parser executor.CommandParser) error {
					var cmd model.StopBlinkingCommand
					parser(&cmd)
					return deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
				}),
			),
			hubman.WithOnConfigRefresh(func(configuration core.AgentConfiguration) {
				update, _ := configuration.User.(*config.UserConfig)
				deviceManager.UpdateDevices(update.MidiDevices)
			}),
			hubman.WithCheckRegistry(checkManager),
		),
	)

	deviceManager.UpdateDevices(userConfig.MidiDevices)

	<-app.WaitShutdown()
}
