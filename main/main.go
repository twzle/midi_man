package main

import (
	"flag"
	"git.miem.hse.ru/hubman/hubman-lib"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"git.miem.hse.ru/hubman/hubman-lib/executor"
	"gitlab.com/gomidi/midi/v2"
	"log"
	"midi_manipulator/pkg/commands"
	"midi_manipulator/pkg/config"
	midiExecutor "midi_manipulator/pkg/executor"
	midiManipulator "midi_manipulator/pkg/manipulator"
	midiSignals "midi_manipulator/pkg/signals"
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

	midiExecutorInstance := midiExecutor.MidiExecutor{}
	go midiExecutorInstance.Run(cfg.MidiConfig)

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
			hubman.WithCommand(commands.TurnLightOnCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd commands.TurnLightOnCommand
					parser(&cmd)
					midiExecutorInstance.TurnLightOnHandler(cmd)
				}),
			hubman.WithCommand(commands.TurnLightOffCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd commands.TurnLightOffCommand
					parser(&cmd)
					midiExecutorInstance.TurnLightOffHandler(cmd)
				}),
			hubman.WithCommand(commands.SingleBlinkCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd commands.SingleBlinkCommand
					parser(&cmd)
				}),
			hubman.WithCommand(commands.SingleReversedBlinkCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd commands.SingleReversedBlinkCommand
					parser(&cmd)
				}),
			hubman.WithCommand(commands.ContinuousBlinkCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd commands.ContinuousBlinkCommand
					parser(&cmd)
				})),
	)
	shutdown := app.WaitShutdown()

	midiManipulatorInstance := midiManipulator.MidiManipulator{}
	go midiManipulatorInstance.Run(cfg.MidiConfig, signals, shutdown)

	<-app.WaitShutdown()
}
