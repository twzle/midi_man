package executor

import (
	"gitlab.com/gomidi/midi/v2"
	"midi_manipulator/pkg/config"
	"slices"
)

const (
	NoteOffCode       = 128
	NoteOnCode        = 144
	ControlChangeCode = 176
)

func (me *MidiExecutor) getTurnLightOffMessage(keyCode int) midi.Message {
	idx := slices.IndexFunc(me.backlightConfig.Backlight, func(c config.Backlight) bool { return c.Key == keyCode })

	if idx == -1 {
		return nil
	}

	var keyBacklight = me.backlightConfig.Backlight[idx].Statuses.Off

	var msg midi.Message
	if keyBacklight.MidiType == "NoteOff" {
		var msg_type = NoteOffCode + keyBacklight.MidiChannel
		var key = keyCode
		var velocity = keyBacklight.MidiVelocity
		msg = midi.Message{byte(msg_type), byte(key), byte(velocity)}
	} else if keyBacklight.MidiType == "ControlChange" {
		var msg_type = ControlChangeCode + keyBacklight.MidiChannel
		var key = keyCode
		var velocity = keyBacklight.MidiVelocity
		msg = midi.Message{byte(msg_type), byte(key), byte(velocity)}
	}

	return msg
}

func (me *MidiExecutor) getTurnLightOnMessage(keyCode int) midi.Message {
	idx := slices.IndexFunc(me.backlightConfig.Backlight, func(c config.Backlight) bool { return c.Key == keyCode })

	if idx == -1 {
		return nil
	}

	var keyBacklight = me.backlightConfig.Backlight[idx].Statuses.On

	var msg midi.Message
	if keyBacklight.MidiType == "NoteOn" {
		var msg_type = NoteOnCode + keyBacklight.MidiChannel
		var key = keyCode
		var velocity = keyBacklight.MidiVelocity
		msg = midi.Message{byte(msg_type), byte(key), byte(velocity)}
	} else if keyBacklight.MidiType == "ControlChange" {
		var msg_type = ControlChangeCode + keyBacklight.MidiChannel
		var key = keyCode
		var velocity = keyBacklight.MidiVelocity
		msg = midi.Message{byte(msg_type), byte(key), byte(velocity)}
	}

	return msg
}
