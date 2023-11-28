package executor

import (
	"gitlab.com/gomidi/midi/v2"
)

const (
	NoteOffCode       = 128
	NoteOnCode        = 144
	ControlChangeCode = 176
)

func GetTurnLightOffMessage(keyCode int) midi.Message {
	msg := midi.Message{byte(NoteOffCode + 1), byte(keyCode), byte(0)}
	return msg
}

func GetTurnLightOnMessage(keyCode int) midi.Message {
	msg := midi.Message{byte(NoteOnCode + 1), byte(keyCode), byte(2)}
	return msg
}
