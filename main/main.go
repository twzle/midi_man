package main

import (
	"flag"
	"git.miem.hse.ru/hubman/hubman-lib"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"git.miem.hse.ru/hubman/hubman-lib/executor"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	"log"
	"midi_manipulator/pkg/config"
	core2 "midi_manipulator/pkg/core"
	"midi_manipulator/pkg/utils"
	"net"
)

func main() {
	defer midi.CloseDriver()

	configPath := flag.String("conf_path", "configs/config.json", "set configs path")
	flag.Parse()

	cfg, err := config.InitConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	setupApp(cfg)
}

func setupApp(cfg *config.Config) {
	deviceManager := core2.NewDeviceManager()
	defer deviceManager.Close()

	agentConf := core.AgentConfiguration{
		System: &core.SystemConfig{
			Server: &core.InterfaceConfig{
				IP:   net.ParseIP(cfg.AppConfig.IPAddr),
				Port: cfg.AppConfig.Port,
			},
			RedisUrl: cfg.RedisConfig.URL,
		},
		User:            cfg.MidiConfig,
		ParseUserConfig: func(data []byte) (core.Configuration, error) { return config.ParseConfigFromBytes(data) },
	}

	signals := deviceManager.GetSignals()
	app := hubman.NewAgentApp(
		agentConf,
		hubman.WithManipulator(
			hubman.WithSignal[utils.NotePushed](),
			hubman.WithSignal[utils.NoteHold](),
			hubman.WithSignal[utils.NoteReleased](),
			hubman.WithSignal[utils.ControlPushed](),
			hubman.WithChannel(signals),
		),
		hubman.WithExecutor(
			hubman.WithCommand(utils.TurnLightOnCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd utils.TurnLightOnCommand
					parser(&cmd)
					err := deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					if err != nil {
						log.Println(err)
					}
				}),
			hubman.WithCommand(utils.TurnLightOffCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd utils.TurnLightOffCommand
					parser(&cmd)
					err := deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					if err != nil {
						log.Println(err)
					}
				}),
			hubman.WithCommand(utils.SingleBlinkCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd utils.SingleBlinkCommand
					parser(&cmd)
					err := deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					if err != nil {
						log.Println(err)
					}
				}),
			hubman.WithCommand(utils.SingleReversedBlinkCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd utils.SingleReversedBlinkCommand
					parser(&cmd)
					err := deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					if err != nil {
						log.Println(err)
					}
				}),
			hubman.WithCommand(utils.ContinuousBlinkCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd utils.ContinuousBlinkCommand
					parser(&cmd)
					err := deviceManager.ExecuteOnDevice(cmd.DeviceAlias, cmd)
					if err != nil {
						log.Println(err)
					}
				})),
		hubman.WithOnConfigRefresh(func(configuration core.AgentConfiguration) {
			update, _ := configuration.User.([]config.MidiConfig)
			deviceManager.UpdateDevices(update)
		}),
	)
	shutdown := app.WaitShutdown()

	deviceManager.SetShutdownChannel(shutdown)
	go deviceManager.UpdateDevices(cfg.MidiConfig)

	<-app.WaitShutdown()
}
