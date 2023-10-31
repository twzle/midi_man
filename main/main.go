package main

import (
	"flag"
	"git.miem.hse.ru/hubman/hubman-lib"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"log"
	"midi_manipulator/pkg/config"
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
	)
	shutdown := app.WaitShutdown()

	midiManipulator := manipulator.MidiManipulator{}
	midiManipulator.Run(cfg.MIDIConfig, signals, shutdown)

}
