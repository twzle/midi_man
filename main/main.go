package main

import (
	"flag"
	"git.miem.hse.ru/hubman/hubman-lib"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"git.miem.hse.ru/hubman/hubman-lib/executor"
	"log"
	"midi_manipulator/pkg/commands"
	"midi_manipulator/pkg/config"
	midiExecutor "midi_manipulator/pkg/executor"
	midiManipulator "midi_manipulator/pkg/manipulator"
	midiSignals "midi_manipulator/pkg/signals"
	"net"
)

func main() {
	configPath := flag.String("conf_path", "configs/config.yaml", "set configs path")
	flag.Parse()

	cfg, err := config.InitConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	setupManipulatorApp(*cfg)
	setupExecutorApp(*cfg)
}

func setupManipulatorApp(cfg config.Config) {
	manipulatorAgentConf := core.AgentConfiguration{
		System: &core.SystemConfig{
			Server: &core.InterfaceConfig{
				IP:   net.ParseIP(cfg.ManipulatorConfig.IPAddr),
				Port: cfg.ManipulatorConfig.Port,
			},
			RedisUrl: cfg.RedisConfig.URL,
		},
		User:            cfg.MIDIConfig,
		ParseUserConfig: func(data []byte) (core.Configuration, error) { return config.ParseConfigFromBytes(data) },
	}

	signals := make(chan core.Signal)
	manipulatorApp := hubman.NewAgentApp(
		manipulatorAgentConf,
		hubman.WithManipulator(
			hubman.WithSignal[midiSignals.NotePushed](),
			hubman.WithSignal[midiSignals.NoteHold](),
			hubman.WithSignal[midiSignals.NoteReleased](),
			hubman.WithSignal[midiSignals.ControlPushed](),
			hubman.WithChannel(signals),
		),
	)
	shutdown := manipulatorApp.WaitShutdown()

	midiManipulatorInstance := midiManipulator.MidiManipulator{}
	go midiManipulatorInstance.Run(cfg.MIDIConfig, signals, shutdown)
}

func setupExecutorApp(cfg config.Config) {
	executorAgentConf := core.AgentConfiguration{
		System: &core.SystemConfig{
			Server: &core.InterfaceConfig{
				IP:   net.ParseIP(cfg.ExecutorConfig.IPAddr),
				Port: cfg.ExecutorConfig.Port,
			},
			RedisUrl: cfg.RedisConfig.URL,
		},
		User:            cfg.MIDIConfig,
		ParseUserConfig: func(data []byte) (core.Configuration, error) { return config.ParseConfigFromBytes(data) },
	}

	midiExecutorInstance := midiExecutor.MidiExecutor{}
	go midiExecutorInstance.StartupIllumination(cfg.MIDIConfig)

	executorApp := hubman.NewAgentApp(
		executorAgentConf,
		hubman.WithExecutor(hubman.WithCommand(commands.TurnLightOnCommand{},
			func(command core.SerializedCommand, parser executor.CommandParser) {
				var cmd commands.TurnLightOnCommand
				parser(&cmd)
				midiExecutorInstance.TurnLightOn(cmd, cfg.MIDIConfig)
			})),
		hubman.WithExecutor(hubman.WithCommand(commands.TurnLightOffCommand{},
			func(command core.SerializedCommand, parser executor.CommandParser) {
				var cmd commands.TurnLightOffCommand
				parser(&cmd)
				midiExecutorInstance.TurnLightOff(cmd, cfg.MIDIConfig)
			})),
		hubman.WithExecutor(hubman.WithCommand(commands.SingleBlinkCommand{},
			func(command core.SerializedCommand, parser executor.CommandParser) {
				var cmd commands.SingleBlinkCommand
				parser(&cmd)
			})),
		hubman.WithExecutor(hubman.WithCommand(commands.SingleReversedBlinkCommand{},
			func(command core.SerializedCommand, parser executor.CommandParser) {
				var cmd commands.SingleReversedBlinkCommand
				parser(&cmd)
			})),
		hubman.WithExecutor(hubman.WithCommand(commands.ContinuousBlinkCommand{},
			func(command core.SerializedCommand, parser executor.CommandParser) {
				var cmd commands.SingleReversedBlinkCommand
				parser(&cmd)
			})),
	)

	<-executorApp.WaitShutdown()
}
