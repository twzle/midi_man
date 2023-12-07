package core

import (
	"fmt"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	_ "gitlab.com/gomidi/midi/v2"
	"io"
	"midi_manipulator/pkg/config"
	"midi_manipulator/pkg/utils"
	"sync"
)

type DeviceManager struct {
	io.Closer
	devices map[string]*MidiDevice
	mutex   sync.Mutex
}

func (dm *DeviceManager) AddDevice(device *MidiDevice) error {
	if _, found := dm.devices[device.GetAlias()]; found {
		return fmt.Errorf("device {%s} already exists", device.GetAlias())
	}

	dm.devices[device.GetAlias()] = device
	return nil
}

func (dm *DeviceManager) RemoveDevice(alias string) error {
	if _, found := dm.devices[alias]; !found {
		return fmt.Errorf("device {%s} doesn't exist", alias)
	}

	delete(dm.devices, alias)
	return nil
}

func (dm *DeviceManager) SetActiveDevice(alias string) error {
	if _, found := dm.devices[alias]; !found {
		return fmt.Errorf("device {%s} doesn't exist", alias)
	}

	if dm.devices[alias].active {
		return fmt.Errorf("device {%s} is already active", alias)
	}

	dm.devices[alias].SetActive()
	return nil
}

func (dm *DeviceManager) ExecuteOnDevice(alias string, command utils.MidiCommand) error {
	if _, found := dm.devices[alias]; !found {
		return fmt.Errorf("device {%s} doesn't exist", alias)
	}

	if dm.devices[alias].active {
		return fmt.Errorf("device {%s} is already active", alias)
	}

	err := dm.devices[alias].ExecuteCommand(command)

	if err != nil {
		return err
	}

	return nil
}

func (dm *DeviceManager) Run(config []config.MidiConfig, signals chan<- core.Signal, shutdown <-chan bool) {
	dm.devices = make(map[string]*MidiDevice)

	for _, deviceConfig := range config {
		midiDevice := MidiDevice{}
		midiDevice.applyConfiguration(deviceConfig)

		err := midiDevice.connectDevice()
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = dm.AddDevice(&midiDevice)
		if err != nil {
			return
		}

		fmt.Printf("MIDI DEVICE {%s} RUNNING ...\n", midiDevice.name)
		midiDevice.startupIllumination()
		go midiDevice.listen(signals, shutdown)
	}
}
