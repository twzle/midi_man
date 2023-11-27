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
	device MidiDevice
	mutex  sync.Mutex
}

type MidiDevice struct {
	name  string
	ports MidiPorts
}

type MidiPorts struct {
	out *drivers.Out
}

func (me *MidiExecutor) TurnLightOn(cmd commands.TurnLightOnCommand) {
	msg := me.getTurnLightOnMessage(cmd.KeyCode)
	if msg != nil && me.device.ports.out != nil {
		(*me.device.ports.out).Send(msg)
	}
}

func (me *MidiExecutor) TurnLightOff(cmd commands.TurnLightOffCommand) {
	msg := me.getTurnLightOffMessage(cmd.KeyCode)
	if msg != nil && me.device.ports.out != nil {
		(*me.device.ports.out).Send(msg)
	}
}

func (me *MidiExecutor) startupIllumination(config config.MIDIConfig) {
	if config.DeviceName == "MPD226" {
		// AKAI MPD226 DIV'S
		for i := 0; i < 4; i++ {
			msg := me.getTurnLightOnMessage(i)
			time.Sleep(time.Millisecond * 50)
			(*me.device.ports.out).Send(msg)
		}

		for i := 0; i < 4; i++ {
			msg := me.getTurnLightOffMessage(i)
			time.Sleep(time.Millisecond * 50)
			(*me.device.ports.out).Send(msg)
		}

		// AKAI MPD226 PADS
		for i := 60; i < 88; i++ {
			msg := me.getTurnLightOnMessage(i)
			time.Sleep(time.Millisecond * 50)
			(*me.device.ports.out).Send(msg)
		}

		for i := 60; i < 88; i++ {
			msg := me.getTurnLightOffMessage(i)
			time.Sleep(time.Millisecond * 50)
			(*me.device.ports.out).Send(msg)
		}
	}
}

func (me *MidiExecutor) initializeDevice(config config.MIDIConfig) error {
	outPort, err := me.getPortsByDeviceName(config.DeviceName)

	if outPort == nil || err != nil {
		return err
	}

	err = me.connectDevice(outPort)

	if err != nil {
		return err
	}

	return nil
}
func (me *MidiExecutor) getPortsByDeviceName(deviceName string) (drivers.Out, error) {
	outPort, err := midi.FindOutPort(deviceName)

	if err != nil {
		log.Println("Output port was not found")
		return nil, err
	}

	return outPort, nil
}

func (me *MidiExecutor) connectDevice(outPort drivers.Out) error {
	port, err := midi.OutPort(outPort.Number())
	if err != nil {
		return err
	}

	me.device.ports.out = &port
	return nil
}

func (me *MidiExecutor) Run(config config.MIDIConfig) {
	err := me.initializeDevice(config)

	if err != nil {
		return
	}

	me.startupIllumination(config)
}
