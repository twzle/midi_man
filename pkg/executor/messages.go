package executor

import (
	"gitlab.com/gomidi/midi/v2"
)

const (
	NoteOffCode       = 128
	NoteOnCode        = 144
	ControlChangeCode = 176
)

func (me *MidiExecutor) getTurnLightOffMessage(keyCode int) midi.Message {
	msg := midi.Message{byte(NoteOffCode), byte(keyCode), byte(0)}
	return msg
}

func (me *MidiExecutor) getTurnLightOnMessage(keyCode int) midi.Message {
	msg := midi.Message{byte(NoteOnCode), byte(keyCode), byte(0)}
	return msg
}
