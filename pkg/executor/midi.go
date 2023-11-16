package executor

import (
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"log"
	"midi_manipulator/pkg/commands"
	"midi_manipulator/pkg/config"
	"sync"
	"time"
)

type MidiExecutor struct {
	device       MidiDevice
	backlightCtx BacklightContext
	mutex        sync.Mutex
}

type MidiDevice struct {
	name  string
	ports MidiPorts
}

type MidiPorts struct {
	out *drivers.Out
}

func (me *MidiExecutor) StartupIllumination(config config.MIDIConfig) {
	outPort := me.getPortsByDeviceName(config.DeviceName)

	if outPort == nil {
		return
	}

	me.connectDevice(outPort)

	me.illuminate(config)
}

func (me *MidiExecutor) TurnLightOn(cmd commands.TurnLightOnCommand, config config.MIDIConfig) {
	outPort := me.getPortsByDeviceName(config.DeviceName)

	if outPort == nil {
		return
	}

	me.connectDevice(outPort)

	msg := me.getTurnLightOnMessage(cmd.KeyCode)
	if msg != nil {
		(*me.device.ports.out).Send(msg)
	}
}

func (me *MidiExecutor) TurnLightOff(cmd commands.TurnLightOffCommand, config config.MIDIConfig) {
	outPort := me.getPortsByDeviceName(config.DeviceName)

	if outPort == nil {
		return
	}

	me.connectDevice(outPort)

	msg := me.getTurnLightOffMessage(cmd.KeyCode)
	if msg != nil {
		(*me.device.ports.out).Send(msg)
	}
}

func (me *MidiExecutor) getTurnLightOnMessage(keyCode int) midi.Message {
	var msg midi.Message
	if keyCode >= 59 && keyCode <= 87 {
		msg = midi.Message{145, byte(keyCode), 2}
	} else if keyCode >= 0 && keyCode <= 3 {
		msg = midi.Message{177, byte(keyCode), 127}
	}
	return msg
}

func (me *MidiExecutor) getTurnLightOffMessage(keyCode int) midi.Message {
	var msg midi.Message
	if keyCode >= 59 && keyCode <= 87 {
		msg = midi.Message{129, byte(keyCode), 2}
	} else if keyCode >= 0 && keyCode <= 3 {
		msg = midi.Message{177, byte(keyCode), 0}
	}
	return msg
}

func (me *MidiExecutor) illuminate(config config.MIDIConfig) {
	// AKAI MPD226 DIV'S
	for i := 0; i < 4; i++ {
		msg := midi.Message{177, byte(i), 127}
		time.Sleep(time.Millisecond * 50)
		(*me.device.ports.out).Send(msg)
	}

	for i := 0; i < 4; i++ {
		msg := midi.Message{177, byte(i), 0}
		time.Sleep(time.Millisecond * 50)
		(*me.device.ports.out).Send(msg)
	}

	// AKAI MPD226 PADS
	for i := 60; i < 88; i++ {
		msg := midi.Message{145, byte(i), 4}
		time.Sleep(time.Millisecond * 50)
		(*me.device.ports.out).Send(msg)
	}

	for i := 60; i < 88; i++ {
		msg := midi.Message{129, byte(i), 2}
		time.Sleep(time.Millisecond * 50)
		(*me.device.ports.out).Send(msg)
	}
}

func (me *MidiExecutor) getPortsByDeviceName(deviceName string) drivers.Out {
	outPort, err := midi.FindOutPort(deviceName)

	if err != nil {
		log.Println("Output port was not found")
		return nil
	}

	return outPort
}

func (me *MidiExecutor) connectDevice(outPort drivers.Out) {
	port, err := midi.OutPort(outPort.Number())
	if err != nil {
		return
	}

	me.device.ports.out = &port
}
