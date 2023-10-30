package main

import (
	"git.miem.hse.ru/hubman/hubman-lib"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"net"
)

func main() {
	agentConf := core.AgentConfiguration{
		System: &core.SystemConfig{
			Server: &core.InterfaceConfig{
				IP:   net.IPv4(127, 0, 0, 1),
				Port: 8095,
			},
			RedisUrl: "redis://127.0.0.1:6379",
		},
	}

	signals := make(chan core.Signal)
	app := hubman.NewAgentApp(
		agentConf,
		hubman.WithManipulator(
			hubman.WithSignal[NotePushed](),
			hubman.WithSignal[NoteHold](),
			hubman.WithSignal[NoteReleased](),
			hubman.WithSignal[ControlPushed](),
			hubman.WithChannel(signals),
		),
	)

	shutdown := app.WaitShutdown()
	Run(signals, shutdown)

}
