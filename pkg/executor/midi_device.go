package executor

import (
	"fmt"
	"gitlab.com/gomidi/midi/v2/drivers"
	"midi_manipulator/pkg/commands"
	"midi_manipulator/pkg/config"
	"time"
)

type MidiDevice struct {
	name string
	port MidiPort
}

type MidiPort struct {
	out *drivers.Out
}

func (md *MidiDevice) setDeviceName(deviceName string) {
	md.name = deviceName
}

func (md *MidiDevice) applyConfiguration(config config.MidiConfig) {
	md.setDeviceName(config.DeviceName)
}

func (md *MidiDevice) turnLightOn(cmd commands.TurnLightOnCommand) {
	msg := GetTurnLightOnMessage(cmd.KeyCode)
	if msg != nil && md.port.out != nil {
		(*md.port.out).Send(msg)
	}
}

func (md *MidiDevice) turnLightOff(cmd commands.TurnLightOffCommand) {
	msg := GetTurnLightOffMessage(cmd.KeyCode)
	if msg != nil && md.port.out != nil {
		(*md.port.out).Send(msg)
	}
}

func (md *MidiDevice) startupIllumination() {
	if md.name == "MPD226" {
		// AKAI MPD226 DIV'S
		for i := 0; i < 4; i++ {
			msg := GetTurnLightOnMessage(i)
			time.Sleep(time.Millisecond * 50)
			(*md.port.out).Send(msg)
		}

		for i := 0; i < 4; i++ {
			msg := GetTurnLightOffMessage(i)
			time.Sleep(time.Millisecond * 50)
			(*md.port.out).Send(msg)
		}

		// AKAI MPD226 PADS
		for i := 60; i < 88; i++ {
			msg := GetTurnLightOnMessage(i)
			time.Sleep(time.Millisecond * 50)
			err := (*md.port.out).Send(msg)

			if err != nil {
				fmt.Println(err)
			}
		}

		for i := 60; i < 88; i++ {
			msg := GetTurnLightOffMessage(i)
			time.Sleep(time.Millisecond * 50)
			(*md.port.out).Send(msg)
		}
	}
}
