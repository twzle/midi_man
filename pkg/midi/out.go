package midi

import (
	"midi_manipulator/pkg/backlight"
	"midi_manipulator/pkg/model"
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

func (md *MidiDevice) startupIllumination(config *backlight.Decoded_DeviceBacklightConfig) {
	for i := 0; i < 4; i++ {
		sequence := backlight.TurnLightOn(config, md.name, i, "red")
		if len(sequence) == 0 {
			continue
		}

		(*md.ports.out).Send(sequence)
	}

	for i := 0; i < 4; i++ {
		sequence := backlight.TurnLightOff(config, md.name, i, "red")
		if len(sequence) == 0 {
			continue
		}

		(*md.ports.out).Send(sequence)
	}

	// AKAI MPD226 PADS
	for i := 60; i < 88; i++ {
		sequence := backlight.TurnLightOn(config, md.name, i, "red")
		if len(sequence) == 0 {
			continue
		}

		(*md.ports.out).Send(sequence)
	}

	for i := 60; i < 88; i++ {
		sequence := backlight.TurnLightOff(config, md.name, i, "red")
		if len(sequence) == 0 {
			continue
		}

		(*md.ports.out).Send(sequence)
	}
}
