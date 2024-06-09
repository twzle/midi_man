package model

import (
	"gitlab.com/gomidi/midi/v2"
)

const (
	NoteOffCode       = 128
	NoteOnCode        = 144
	ControlChangeCode = 176
)

// Function creating MIDi-compatible message from keycode to turn light off for single component
func GetTurnLightOffMessage(keyCode int) midi.Message {
	msg := midi.Message{byte(NoteOffCode), byte(keyCode), byte(0)}
	return msg
}

// Function creating MIDi-compatible message from keycode to turn light on for single component
func GetTurnLightOnMessage(keyCode int) midi.Message {
	msg := midi.Message{byte(NoteOnCode), byte(keyCode), byte(127)}
	return msg
}
