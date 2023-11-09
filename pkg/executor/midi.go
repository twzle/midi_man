package executor

import (
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"log"
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
	in  drivers.In
	out drivers.Out
}

func (me *MidiExecutor) startupIllumination() {
	// AKAI MPD226 DIV'S
	for i := 0; i < 4; i++ {
		msg := midi.Message{177, byte(i), 127}
		time.Sleep(time.Millisecond * 50)
		me.device.ports.out.Send(msg)
	}

	for i := 0; i < 4; i++ {
		msg := midi.Message{177, byte(i), 0}
		time.Sleep(time.Millisecond * 50)
		me.device.ports.out.Send(msg)
	}

	// AKAI MPD226 PADS
	for i := 60; i < 88; i++ {
		msg := midi.Message{145, byte(i), 4}
		time.Sleep(time.Millisecond * 50)
		me.device.ports.out.Send(msg)
	}

	for i := 60; i < 88; i++ {
		msg := midi.Message{129, byte(i), 2}
		time.Sleep(time.Millisecond * 50)
		me.device.ports.out.Send(msg)
	}
}

func (me *MidiExecutor) getPortsByDeviceName(deviceName string) (drivers.In, drivers.Out) {
	inPort, err := midi.FindInPort(deviceName)
	if err != nil {
		log.Println("Input port was not found")
		return nil, nil
	}

	outPort, err := midi.FindOutPort(deviceName)
	if err != nil {
		log.Println("Output port was not found")
		return nil, nil
	}

	return inPort, outPort
}

func (me *MidiExecutor) connectDevice(inPort drivers.In, outPort drivers.Out) {
	me.device.ports.in, _ = midi.InPort(inPort.Number())
	me.device.ports.out, _ = midi.OutPort(outPort.Number())
}

func (me *MidiExecutor) Run(config config.MIDIConfig) {
	defer midi.CloseDriver()
	inPort, outPort := me.getPortsByDeviceName(config.DeviceName)
	if inPort == nil || outPort == nil {
		return
	}

	me.connectDevice(inPort, outPort)
	defer me.device.ports.in.Close()
	defer me.device.ports.out.Close()
	me.startupIllumination()
}
