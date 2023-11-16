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
		User:            cfg.MIDIConfig,
		ParseUserConfig: func(data []byte) (core.Configuration, error) { return config.ParseConfigFromBytes(data) },
	}

	midiExecutorInstance := midiExecutor.MidiExecutor{}
	go midiExecutorInstance.StartupIllumination(cfg.MIDIConfig)

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
					midiExecutorInstance.TurnLightOn(cmd, cfg.MIDIConfig)
				}),
			hubman.WithCommand(commands.TurnLightOffCommand{},
				func(command core.SerializedCommand, parser executor.CommandParser) {
					var cmd commands.TurnLightOffCommand
					parser(&cmd)
					midiExecutorInstance.TurnLightOff(cmd, cfg.MIDIConfig)
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
	go midiManipulatorInstance.Run(cfg.MIDIConfig, signals, shutdown)

	<-app.WaitShutdown()

	return
}
