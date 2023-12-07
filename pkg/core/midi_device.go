package core

import (
	"fmt"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"log"
	"midi_manipulator/pkg/config"
	"midi_manipulator/pkg/utils"
	"sync"
	"time"
)

type MidiDevice struct {
	name        string
	active      bool
	ports       MidiPorts
	clickBuffer ClickBuffer
	holdDelta   time.Duration
	mutex       sync.Mutex
}

type MidiPorts struct {
	in  *drivers.In
	out *drivers.Out
}

func (md *MidiDevice) GetAlias() string {
	return md.name
}

func (md *MidiDevice) SetActive() {
	md.active = true
}

func (md *MidiDevice) ExecuteCommand(command utils.MidiCommand) error {
	return nil
}

func (md *MidiDevice) connectDevice() error {
	var err error
	in_err := md.connectInPort()
	out_err := md.connectOutPort()

	if in_err != nil || out_err != nil {
		err = fmt.Errorf("Connection of device \"{%s}\" failed\n", md.name)
	}
	return err
}

func (md *MidiDevice) connectOutPort() error {
	fmt.Println(midi.GetOutPorts())
	port, err := midi.FindOutPort(md.name)
	if err != nil {
		log.Printf("Output port named {%s} was not found\n", md.name)
		return err
	}

	port, err = midi.OutPort(port.Number())
	if err != nil {
		log.Printf("Output port #{%d} was not found\n", port.Number())
		return err
	}

	md.ports.out = &port
	return nil
}

func (md *MidiDevice) connectInPort() error {
	fmt.Println(midi.GetInPorts())
	port, err := midi.FindInPort(md.name)
	if err != nil {
		log.Printf("Input port named {%s} was not found", md.name)
		return err
	}

	port, err = midi.InPort(port.Number())
	if err != nil {
		log.Printf("Input port #{%d} was not found\n", port.Number())
		return err
	}

	md.ports.in = &port
	return nil
}

func (md *MidiDevice) applyConfiguration(config config.MidiConfig) {
	md.name = config.DeviceName
	md.holdDelta = time.Duration(float64(time.Second) * config.HoldDelta)
	md.clickBuffer = make(map[uint8]*KeyContext)
}
