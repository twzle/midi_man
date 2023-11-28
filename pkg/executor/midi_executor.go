package executor

import (
	"fmt"
	"gitlab.com/gomidi/midi/v2"
	"midi_manipulator/pkg/commands"
	"midi_manipulator/pkg/config"
	"sync"
)

type MidiExecutor struct {
	devices []*MidiDevice
	mutex   sync.Mutex
}

func (me *MidiExecutor) TurnLightOnHandler(cmd commands.TurnLightOnCommand) {
	for _, device := range me.devices {
		device.turnLightOn(cmd)
	}
}

func (me *MidiExecutor) TurnLightOffHandler(cmd commands.TurnLightOffCommand) {
	for _, device := range me.devices {
		device.turnLightOff(cmd)
	}
}

func (me *MidiExecutor) addDevice(device *MidiDevice) {
	me.devices = append(me.devices, device)
}

func (me *MidiExecutor) removeDevice(device *MidiDevice) {
	for idx, x := range me.devices {
		if device == x {
			me.devices = append(me.devices[:idx], me.devices[idx+1:]...)
			fmt.Println("Device deleted: ", device.name)
		}
	}
}

func (me *MidiExecutor) connectDevice(device *MidiDevice) error {
	port, err := midi.FindOutPort(device.name)
	if err != nil {
		fmt.Printf("Output port {%s} not found\n", device.name)
		return err
	}

	port, err = midi.OutPort(port.Number())
	if err != nil {
		fmt.Printf("Output port #{%d} not found\n", port.Number())
		return err
	}

	device.port.out = &port
	return nil
}

func (me *MidiExecutor) Run(config []config.MidiConfig) {
	for _, device := range config {
		midiDevice := MidiDevice{}
		midiDevice.applyConfiguration(device)

		err := me.connectDevice(&midiDevice)

		if err != nil {
			continue
		}

		me.addDevice(&midiDevice)

		fmt.Printf("WRITING TO DEVICE {%s} ...\n", device.DeviceName)
		midiDevice.startupIllumination()
	}
}
