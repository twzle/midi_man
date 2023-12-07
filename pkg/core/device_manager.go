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
	devices  map[string]*MidiDevice
	mutex    sync.Mutex
	shutdown <-chan bool
	signals  chan<- core.Signal
}

func (dm *DeviceManager) AddDevice(device *MidiDevice) error {
	_, found := dm.devices[device.GetAlias()]

	if found {
		return fmt.Errorf("device {%s} already exists", device.GetAlias())
	}

	dm.devices[device.GetAlias()] = device

	dm.runDevice(device)
	return nil
}

func (dm *DeviceManager) RemoveDevice(alias string) error {
	device, found := dm.devices[alias]
	if !found {
		return fmt.Errorf("device {%s} doesn't exist", device.GetAlias())
	}

	err := device.StopDevice()
	if err != nil {
		return err
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

func (dm *DeviceManager) UpdateDevices(midiConfig []config.MidiConfig) {
	for _, deviceConfig := range midiConfig {
		if _, found := dm.devices[deviceConfig.DeviceName]; !found {
			device, err := dm.initializeDevice(deviceConfig)

			if err != nil {
				return
			}

			err = dm.AddDevice(device)

			if err != nil {
				return
			}

			dm.runDevice(device)
		}
	}

	for _, device := range dm.devices {
		var found = false

		for _, v := range midiConfig {
			if v.DeviceName == device.GetAlias() {
				found = true
				break
			}
		}

		if !found {
			dm.stopDevice(device)
			err := dm.RemoveDevice(device.GetAlias())
			if err != nil {
				return
			}
		}

	}
}

func (dm *DeviceManager) initializeDevice(deviceConfig config.MidiConfig) (*MidiDevice, error) {
	midiDevice := MidiDevice{}
	midiDevice.applyConfiguration(deviceConfig)

	err := midiDevice.connectDevice()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &midiDevice, nil
}

func (dm *DeviceManager) runDevice(device *MidiDevice) {
	fmt.Printf("MIDI DEVICE {%s} RUNNING ...\n", device.name)
	device.startupIllumination()
	go device.listen(dm.signals, dm.shutdown)
}

func (dm *DeviceManager) stopDevice(device *MidiDevice) {
	fmt.Printf("MIDI DEVICE {%s} STOPPING ...\n", device.name)
	device.stop <- true
}

func (dm *DeviceManager) Run(config []config.MidiConfig, signals chan<- core.Signal, shutdown <-chan bool) {
	dm.devices = make(map[string]*MidiDevice)
	dm.shutdown = shutdown
	dm.signals = signals

	for _, deviceConfig := range config {
		device, err := dm.initializeDevice(deviceConfig)

		if err != nil {
			return
		}

		err = dm.AddDevice(device)

		if err != nil {
			return
		}
	}
}
