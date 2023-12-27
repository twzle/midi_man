package midi

import (
	"midi_manipulator/pkg/backlight"
	"midi_manipulator/pkg/model"
	"time"
)

func (md *MidiDevice) turnLightOn(cmd model.TurnLightOnCommand, backlightConfig *backlight.DecodedDeviceBacklightConfig) {
	msg, _ := backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, "on")
	if msg != nil && md.ports.out != nil {
		(*md.ports.out).Send(msg)
	}
}

func (md *MidiDevice) turnLightOff(cmd model.TurnLightOffCommand, backlightConfig *backlight.DecodedDeviceBacklightConfig) {
	msg, _ := backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, "off")
	if msg != nil && md.ports.out != nil {
		(*md.ports.out).Send(msg)
	}
}

func (md *MidiDevice) startupIllumination(config *backlight.DecodedDeviceBacklightConfig) {
	backlightTimeOffset := time.Duration(config.DeviceBacklightTimeOffset[md.name])
	for _, keyRange := range config.DeviceKeyRangeMap[md.name] {
		for i := keyRange[0]; i <= keyRange[1]; i++ {
			sequence, _ := config.TurnLight(md.name, i, "red", "on")
			if len(sequence) == 0 {
				continue
			}

			time.Sleep(time.Millisecond * backlightTimeOffset)

			(*md.ports.out).Send(sequence)

		}

		for i := keyRange[0]; i <= keyRange[1]; i++ {
			sequence, _ := config.TurnLight(md.name, i, "red", "off")
			if len(sequence) == 0 {
				continue
			}

			time.Sleep(time.Millisecond * backlightTimeOffset)

			(*md.ports.out).Send(sequence)
		}
	}
}
