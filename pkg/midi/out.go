package midi

import (
	"midi_manipulator/pkg/model"
	"time"
)

func (md *MidiDevice) turnLightOn(cmd model.TurnLightOnCommand) {
	msg := model.GetTurnLightOnMessage(cmd.KeyCode)
	if msg != nil && md.ports.out != nil {
		(*md.ports.out).Send(msg)
	}
}

func (md *MidiDevice) turnLightOff(cmd model.TurnLightOffCommand) {
	msg := model.GetTurnLightOffMessage(cmd.KeyCode)
	if msg != nil && md.ports.out != nil {
		(*md.ports.out).Send(msg)
	}
}

func (md *MidiDevice) startupIllumination() {
	if md.name == "MPD226" {
		// AKAI MPD226 DIV'S
		for i := 0; i < 4; i++ {
			msg := model.GetTurnLightOnMessage(i)
			time.Sleep(time.Millisecond * 50)
			(*md.ports.out).Send(msg)
		}

		for i := 0; i < 4; i++ {
			msg := model.GetTurnLightOffMessage(i)
			time.Sleep(time.Millisecond * 50)
			(*md.ports.out).Send(msg)
		}

		// AKAI MPD226 PADS
		for i := 60; i < 88; i++ {
			msg := model.GetTurnLightOnMessage(i)
			time.Sleep(time.Millisecond * 50)
			(*md.ports.out).Send(msg)
		}

		for i := 60; i < 88; i++ {
			msg := model.GetTurnLightOffMessage(i)
			time.Sleep(time.Millisecond * 50)
			(*md.ports.out).Send(msg)
		}
	}
}
