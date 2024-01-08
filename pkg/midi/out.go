package midi

import (
	"midi_manipulator/pkg/backlight"
	"midi_manipulator/pkg/model"
	"time"
)

func (md *MidiDevice) turnLightOn(cmd model.TurnLightOnCommand, backlightConfig *backlight.DecodedDeviceBacklightConfig) {
	msg, _ := backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, backlight.On)
	if msg != nil && md.ports.out != nil {
		(*md.ports.out).Send(msg)
	}
}

func (md *MidiDevice) turnLightOff(cmd model.TurnLightOffCommand, backlightConfig *backlight.DecodedDeviceBacklightConfig) {
	msg, _ := backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, backlight.Off)
	if msg != nil && md.ports.out != nil {
		(*md.ports.out).Send(msg)
	}
}

func (md *MidiDevice) singleBlink(cmd model.SingleBlinkCommand, backlightConfig *backlight.DecodedDeviceBacklightConfig) {
	backlightTimeOffset := time.Duration(backlightConfig.DeviceBacklightTimeOffset[md.name])
	msg, _ := backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, backlight.On)
	if msg != nil && md.ports.out != nil {
		(*md.ports.out).Send(msg)
	}
	time.Sleep(time.Millisecond * backlightTimeOffset)
	msg, _ = backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, backlight.Off)
	if msg != nil && md.ports.out != nil {
		(*md.ports.out).Send(msg)
	}
}

func (md *MidiDevice) singleReversedBlink(cmd model.SingleReversedBlinkCommand, backlightConfig *backlight.DecodedDeviceBacklightConfig) {
	backlightTimeOffset := time.Duration(backlightConfig.DeviceBacklightTimeOffset[md.name])
	msg, _ := backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, backlight.Off)
	if msg != nil && md.ports.out != nil {
		(*md.ports.out).Send(msg)
	}
	time.Sleep(time.Millisecond * backlightTimeOffset)
	msg, _ = backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, backlight.On)
	if msg != nil && md.ports.out != nil {
		(*md.ports.out).Send(msg)
	}
}

func (md *MidiDevice) startupIllumination(config *backlight.DecodedDeviceBacklightConfig) {
	backlightTimeOffset := time.Duration(config.DeviceBacklightTimeOffset[md.name])
	for _, keyRange := range config.DeviceKeyRangeMap[md.name] {
		for i := keyRange[0]; i <= keyRange[1]; i++ {
			sequence, _ := config.TurnLight(md.name, i, "none", backlight.On)
			if len(sequence) == 0 {
				continue
			}

			time.Sleep(time.Millisecond * backlightTimeOffset)
			(*md.ports.out).Send(sequence)
		}

		for i := keyRange[0]; i <= keyRange[1]; i++ {
			sequence, _ := config.TurnLight(md.name, i, "none", backlight.Off)
			if len(sequence) == 0 {
				continue
			}

			time.Sleep(time.Millisecond * backlightTimeOffset)
			(*md.ports.out).Send(sequence)
		}
	}
}
