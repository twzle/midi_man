package main

import (
	"flag"
	"fmt"
	"git.miem.hse.ru/hubman/hubman-lib"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"git.miem.hse.ru/hubman/hubman-lib/executor"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	"log"
	"midi_manipulator/pkg/config"
	core2 "midi_manipulator/pkg/core"
	midiSignals "midi_manipulator/pkg/utils"
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
	deviceManager := core2.DeviceManager{}

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

	signals := make(chan core.Signal)
	app := hubman.NewAgentApp(
		agentConf,
		hubman.WithManipulator(
			hubman.WithSignal[midiSignals.NotePushed](),
			hubman.WithSignal[midiSignals.NoteHold](),
			hubman.WithSignal[midiSignals.NoteReleased](),
			hubman.WithSignal[midiSignals.ControlPushed](),
			hubman.WithChannel(signals),
		),
		hubman.WithExecutor(
			hubman.WithCommand(midiSignals.TurnLightOnCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd midiSignals.TurnLightOnCommand
					parser(&cmd)
					deviceManager.TurnLightOnHandler(cmd)
				}),
			hubman.WithCommand(midiSignals.TurnLightOffCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd midiSignals.TurnLightOffCommand
					parser(&cmd)
					deviceManager.TurnLightOffHandler(cmd)
				}),
			hubman.WithCommand(midiSignals.SingleBlinkCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd midiSignals.SingleBlinkCommand
					parser(&cmd)
				}),
			hubman.WithCommand(midiSignals.SingleReversedBlinkCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd midiSignals.SingleReversedBlinkCommand
					parser(&cmd)
				}),
			hubman.WithCommand(midiSignals.ContinuousBlinkCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd midiSignals.ContinuousBlinkCommand
					parser(&cmd)
				})),
		hubman.WithOnConfigRefresh(func(configuration core.AgentConfiguration) {
			update, ok := configuration.User.([]config.MidiConfig)
			if !ok {
				panic(
					fmt.Sprintf(
						"Refresh config error: expected type %T, received %T",
						config.MidiConfig{},
						configuration.User,
					),
				)
			}
			deviceManager.UpdateDevices(update)
		}),
	)
	shutdown := app.WaitShutdown()

	go deviceManager.Run(cfg.MidiConfig, signals, shutdown)

	<-app.WaitShutdown()
}
