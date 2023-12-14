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
	midiHermophrodite "midi_manipulator/pkg/midi"
	"midi_manipulator/pkg/model"
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
	deviceManager := midiHermophrodite.NewDeviceManager()
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

	deviceManager.UpdateDevices(cfg.MidiConfig)
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
			update, _ := configuration.User.([]config.MidiConfig)
			deviceManager.UpdateDevices(update)
		}),
	)
	<-app.WaitShutdown()
}
