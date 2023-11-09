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
	"midi_manipulator/pkg/manipulator"
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

	agentConf := core.AgentConfiguration{
		System: &core.SystemConfig{
			Server: &core.InterfaceConfig{
				IP:   net.ParseIP(cfg.ManpulatorConfig.IPAddr),
				Port: cfg.ManpulatorConfig.Port,
			},
			RedisUrl: cfg.RedisConfig.URL,
		},
		User:            cfg.MIDIConfig,
		ParseUserConfig: func(data []byte) (core.Configuration, error) { return config.ParseConfigFromBytes(data) },
	}

	signals := make(chan core.Signal)
	manipulator_app := hubman.NewAgentApp(
		agentConf,
		hubman.WithManipulator(
			hubman.WithSignal[midiSignals.NotePushed](),
			hubman.WithSignal[midiSignals.NoteHold](),
			hubman.WithSignal[midiSignals.NoteReleased](),
			hubman.WithSignal[midiSignals.ControlPushed](),
			hubman.WithChannel(signals),
		),
	)
	shutdown := manipulator_app.WaitShutdown()

	midiManipulator := manipulator.MidiManipulator{}
	go midiManipulator.Run(cfg.MIDIConfig, signals, shutdown)

	executor_app := hubman.NewAgentApp(
		agentConf,
		hubman.WithExecutor(hubman.WithCommand(commands.TurnLightOnCommand{},
			func(command core.SerializedCommand, parser executor.CommandParser) {
				var cmd commands.TurnLightOnCommand
				parser(&cmd)
				// midiExecutor.lightOff([]key_codes)
				// midiExecutor.lightOn([]key_codes)
			})),
		hubman.WithExecutor(hubman.WithCommand(commands.TurnLightOffCommand{},
			func(command core.SerializedCommand, parser executor.CommandParser) {
				var cmd commands.TurnLightOffCommand
				parser(&cmd)
				// midiExecutor.lightOff([]key_codes)
			})),
		hubman.WithExecutor(hubman.WithCommand(commands.SingleBlinkCommand{},
			func(command core.SerializedCommand, parser executor.CommandParser) {
				var cmd commands.SingleBlinkCommand
				parser(&cmd)
				// midiExecutor.lightOff([]key_codes)
				// midiExecutor.singleBlink([]key_codes)
			})),
		hubman.WithExecutor(hubman.WithCommand(commands.SingleReversedBlinkCommand{},
			func(command core.SerializedCommand, parser executor.CommandParser) {
				var cmd commands.SingleReversedBlinkCommand
				parser(&cmd)
				// midiExecutor.lightOff([]key_codes)
				// midiExecutor.singleReversedBlink([]key_codes)
			})),
		hubman.WithExecutor(hubman.WithCommand(commands.ContinuousBlinkCommand{},
			func(command core.SerializedCommand, parser executor.CommandParser) {
				var cmd commands.SingleReversedBlinkCommand
				parser(&cmd)
				// midiExecutor.lightOff([]key_codes)
				// midiExecutor.continuousBlink([]key_codes)
			})),
	)

	midiExecutor := midiExecutor.MidiExecutor{}
	midiExecutor.Run(cfg.MIDIConfig)

	<-executor_app.WaitShutdown()

}
