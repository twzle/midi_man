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
	signals  chan core.Signal
}

func (dm *DeviceManager) ExecuteCommand(cmd utils.MidiCommand) {
	for _, device := range dm.devices {
		if device.active {
			err := dm.ExecuteOnDevice(device.GetAlias(), cmd)

			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

func (dm *DeviceManager) AddDevice(device *MidiDevice) error {
	_, found := dm.devices[device.GetAlias()]

	if found {
		return fmt.Errorf("device {%s} already exists", device.GetAlias())
	}

	dm.devices[device.GetAlias()] = device

	err := device.RunDevice(dm.signals, dm.shutdown)

	if err != nil {
		return err
	}

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
	err := dm.devices[alias].ExecuteCommand(command)

	if err != nil {
		return err
	}

	return nil
}

func (dm *DeviceManager) UpdateDevices(midiConfig []config.MidiConfig) {
	var midiConfigMap = make(map[string]config.MidiConfig)

	for _, device := range midiConfig {
		midiConfigMap[device.DeviceName] = device
	}

	for _, deviceConfig := range midiConfigMap {
		if _, found := dm.devices[deviceConfig.DeviceName]; !found {
			device, err := NewDevice(deviceConfig)

			if err != nil {
				continue
			}

			err = dm.AddDevice(device)

			if err != nil {
				continue
			}
		}
	}

	for _, device := range dm.devices {
		if _, found := midiConfigMap[device.GetAlias()]; !found {
			err := dm.RemoveDevice(device.GetAlias())
			if err != nil {
				continue
			}
		}
	}
}

func (dm *DeviceManager) SetShutdownChannel(shutdown <-chan bool) {
	dm.shutdown = shutdown
}

func (dm *DeviceManager) GetSignals() chan core.Signal {
	return dm.signals
}

func NewDeviceManager() *DeviceManager {
	dm := DeviceManager{}
	dm.devices = make(map[string]*MidiDevice)
	dm.signals = make(chan core.Signal)
	dm.shutdown = nil

	return &dm
}
