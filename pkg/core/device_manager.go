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
	closer   io.Closer
	devices  map[string]*MidiDevice
	mutex    sync.Mutex
	shutdown <-chan bool
	signals  chan core.Signal
}

func (dm *DeviceManager) Close() {
	err := dm.closer.Close()

	if err != nil {
		return
	}
}

func (dm *DeviceManager) ExecuteOnDevice(alias string, cmd utils.MidiCommand) error {
	device, found := dm.devices[alias]

	if !found {
		return fmt.Errorf("device with alias {%s} doesn't exist", alias)
	}

	if !device.active {
		return fmt.Errorf("device with alias {%s} is not active", alias)
	}

	err := device.ExecuteCommand(cmd)

	if err != nil {
		return err
	}

	return nil
}

func (dm *DeviceManager) AddDevice(device *MidiDevice) error {
	dm.mutex.Lock()

	_, found := dm.devices[device.GetAlias()]

	if found {
		return fmt.Errorf("device {%s} already exists", device.GetAlias())
	}

	dm.devices[device.GetAlias()] = device

	err := device.RunDevice(dm.signals, dm.shutdown)

	if err != nil {
		return err
	}

	dm.mutex.Unlock()

	return nil
}

func (dm *DeviceManager) RemoveDevice(alias string) error {
	dm.mutex.Lock()

	device, found := dm.devices[alias]
	if !found {
		return fmt.Errorf("device {%s} doesn't exist", device.GetAlias())
	}

	err := device.StopDevice()

	if err != nil {
		return err
	}

	delete(dm.devices, alias)

	dm.mutex.Unlock()

	return nil
}

func (dm *DeviceManager) UpdateDevices(midiConfig []config.MidiConfig) {
	var midiConfigMap = make(map[string]config.MidiConfig)

	for _, device := range midiConfig {
		midiConfigMap[device.DeviceName] = device
	}

	for _, deviceConfig := range midiConfigMap {
		device, found := dm.devices[deviceConfig.DeviceName]
		if !found {
			newDevice, err := NewDevice(deviceConfig)

			if err != nil {
				continue
			}

			err = dm.AddDevice(newDevice)

			if err != nil {
				continue
			}
		} else {
			device.updateConfiguration(deviceConfig)
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
