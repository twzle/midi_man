package main

import (
	"fmt"
	"git.miem.hse.ru/hubman/hubman-lib"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"net"
)

func main() {

	systemConf := &core.SystemConfig{
		Server: &core.InterfaceConfig{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 8094,
		},
		RedisUrl: "redis://127.0.0.1:6379",
	}

	agentConf := &core.AgentConfiguration{
		System: systemConf,
		User:   nil,
		ParseUserConfig: func(bytes []byte) (core.Configuration, error) {
			return nil, nil
		},
	}

	_manipulator := hubman.NewManipulator(*agentConf, hubman.WithSignal[NotePushed](),
		hubman.WithSignal[NoteHold](), hubman.WithSignal[NoteReleased](), hubman.WithSignal[ControlPushed]())
	_ = hubman.NewAgentApp(*agentConf, hubman.WithManipulator(_manipulator))

	signs := make(chan core.Signal)
	go Run(signs)

	for s := range signs {
		if s == nil {
			continue
		}
		err := _manipulator.Process(s)
		fmt.Println(s.Code(), s)
		if err != nil {
			return
		}
	}
}
