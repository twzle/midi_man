package midi

import (
	"midi_manipulator/pkg/backlight"
	"midi_manipulator/pkg/model"
	"time"
)

// Function handles logic of turn light on command
func (md *MidiDevice) turnLightOn(cmd model.TurnLightOnCommand, backlightConfig *backlight.DecodedDeviceBacklightConfig) {
	msg, _ := backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, backlight.On)
	if msg != nil && md.ports.out != nil {
		md.ports.out.Send(msg)
	}
}

// Function handles logic of turn light off command
func (md *MidiDevice) turnLightOff(cmd model.TurnLightOffCommand, backlightConfig *backlight.DecodedDeviceBacklightConfig) {
	msg, _ := backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, backlight.Off)
	if msg != nil && md.ports.out != nil {
		md.ports.out.Send(msg)
	}
}

// Function handles logic of single blink command
func (md *MidiDevice) singleBlink(cmd model.SingleBlinkCommand, backlightConfig *backlight.DecodedDeviceBacklightConfig) {
	backlightTimeOffset := time.Duration(backlightConfig.DeviceBacklightTimeOffset[md.name])
	msg, _ := backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, backlight.On)
	if msg != nil && md.ports.out != nil {
		md.ports.out.Send(msg)
	}
	time.Sleep(time.Millisecond * backlightTimeOffset)
	msg, _ = backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, backlight.Off)
	if msg != nil && md.ports.out != nil {
		md.ports.out.Send(msg)
	}
}

// Function handles logic of single reversed blink command
func (md *MidiDevice) singleReversedBlink(
	cmd model.SingleReversedBlinkCommand,
	backlightConfig *backlight.DecodedDeviceBacklightConfig,
) {
	backlightTimeOffset := time.Duration(backlightConfig.DeviceBacklightTimeOffset[md.name])
	msg, _ := backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, backlight.Off)
	if msg != nil && md.ports.out != nil {
		md.ports.out.Send(msg)
	}
	time.Sleep(time.Millisecond * backlightTimeOffset)
	msg, _ = backlightConfig.TurnLight(md.name, byte(cmd.KeyCode), cmd.ColorName, backlight.On)
	if msg != nil && md.ports.out != nil {
		md.ports.out.Send(msg)
	}
}

// Function handles logic of changing active namespace for device
func (md *MidiDevice) setActiveNamespace(
	cmd model.SetActiveNamespaceCommand,
	_ *backlight.DecodedDeviceBacklightConfig,
) {
	oldNamespace := md.namespace
	md.namespace = cmd.Namespace
	md.sendNamespaceChangedSignal(md.signals, oldNamespace, cmd.Namespace)
}

// Function handles logic of turning light for range of keys of single MIDI-device
func (md *MidiDevice) turnLightKeyRange(
	config *backlight.DecodedDeviceBacklightConfig,
	left, right byte,
	status backlight.StatusName,
	backlightTimeOffset time.Duration,
) {
	for i := left; i <= right; i++ {
		sequence, _ := config.TurnLight(md.name, i, "none", status)
		if len(sequence) == 0 {
			continue
		}

		time.Sleep(time.Millisecond * backlightTimeOffset)
		md.ports.out.Send(sequence)
	}
}

// Function handles startup illumination for MIDI-device
func (md *MidiDevice) startupIllumination(config *backlight.DecodedDeviceBacklightConfig) {
	time.Sleep(md.startupDelay)
	backlightTimeOffset := time.Duration(config.DeviceBacklightTimeOffset[md.name])
	for _, keyRange := range config.DeviceKeyRangeMap[md.name] {
		md.turnLightKeyRange(config, keyRange[0], keyRange[1], backlight.On, backlightTimeOffset)
		md.turnLightKeyRange(config, keyRange[0], keyRange[1], backlight.Off, backlightTimeOffset)
	}
}
