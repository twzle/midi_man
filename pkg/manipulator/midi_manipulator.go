package manipulator

import (
	"fmt"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"log"
	"midi_manipulator/pkg/config"
)

type MidiManipulator struct {
	devices []*MidiDevice
}

func (mm *MidiManipulator) addDevice(device *MidiDevice) {
	mm.devices = append(mm.devices, device)
}

func (mm *MidiManipulator) removeDevice(device *MidiDevice) {
	for idx, x := range mm.devices {
		if device == x {
			mm.devices = append(mm.devices[:idx], mm.devices[idx+1:]...)
			fmt.Println("Device deleted: ", device.name)
		}
	}
}

func (mm *MidiManipulator) connectDevice(midiDevice *MidiDevice) error {
	port, err := midi.FindInPort(midiDevice.name)
	if err != nil {
		log.Printf("Input port named {%s} was not found", midiDevice.name)
		return err
	}

	port, err = midi.InPort(port.Number())
	if err != nil {
		fmt.Printf("Input port #{%d} not found\n", port.Number())
		return err
	}

	midiDevice.ports.in = &port
	return nil
}

func (mm *MidiManipulator) listen(device *MidiDevice, signals chan<- core.Signal, shutdown <-chan bool) {
	stop, err := midi.ListenTo(*device.ports.in, device.getMidiMessage, midi.UseSysEx())

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	for {
		signalSequence := device.messageToSignal()
		select {
		case <-shutdown:
			stop()
			mm.removeDevice(device)
			return
		default:
			for _, signal := range signalSequence {
				device.sendSignal(signals, signal)
			}
		}
	}
}

func (mm *MidiManipulator) Run(config []config.MidiConfig, signals chan<- core.Signal, shutdown <-chan bool) {
	for _, deviceConfig := range config {
		midiDevice := MidiDevice{}
		midiDevice.applyConfiguration(deviceConfig)

		err := mm.connectDevice(&midiDevice)
		if err != nil {
			continue
		}

		mm.addDevice(&midiDevice)

		fmt.Printf("LISTENING INPUT FROM DEVICE {%s} ...\n", midiDevice.name)
		go mm.listen(&midiDevice, signals, shutdown)
	}
}
