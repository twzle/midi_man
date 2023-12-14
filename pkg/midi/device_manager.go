package midi

import (
	"fmt"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	_ "gitlab.com/gomidi/midi/v2"
	"log"
	"midi_manipulator/pkg/config"
	"midi_manipulator/pkg/model"
	"sync"
)

type DeviceManager struct {
	devices map[string]*MidiDevice
	mutex   sync.Mutex
	signals chan core.Signal
}

func (dm *DeviceManager) getDevice(alias string) (*MidiDevice, bool) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	device, ok := dm.devices[alias]
	return device, ok
}

func (dm *DeviceManager) addDevice(device *MidiDevice) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	dm.devices[device.GetAlias()] = device
}

func (dm *DeviceManager) removeDevice(alias string) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	delete(dm.devices, alias)
}

func (dm *DeviceManager) Close() {
	for _, device := range dm.devices {
		err := dm.RemoveDevice(device.GetAlias())
		if err != nil {
			log.Println(err)
		}
	}
}

func (dm *DeviceManager) ExecuteOnDevice(alias string, cmd model.MidiCommand) error {
	device, found := dm.getDevice(alias)

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
	_, found := dm.getDevice(device.GetAlias())

	if found {
		return fmt.Errorf("device {%s} already exists", device.GetAlias())
	}

	dm.addDevice(device)
	err := device.RunDevice(dm.signals)

	if err != nil {
		return err
	}

	return nil
}

func (dm *DeviceManager) RemoveDevice(alias string) error {
	device, found := dm.getDevice(alias)
	if !found {
		return fmt.Errorf("device {%s} doesn't exist", device.GetAlias())
	}

	err := device.StopDevice()

	if err != nil {
		return err
	}

	dm.removeDevice(alias)

	return nil
}

func (dm *DeviceManager) UpdateDevices(midiConfig []config.MidiConfig) {
	var midiConfigMap = make(map[string]config.MidiConfig)
	defer dm.mutex.Unlock()

	dm.mutex.Lock()
	for _, device := range midiConfig {
		midiConfigMap[device.DeviceName] = device
	}
	dm.mutex.Unlock()

	for _, deviceConfig := range midiConfigMap {
		device, found := dm.getDevice(deviceConfig.DeviceName)
		if !found {
			newDevice, err := NewDevice(deviceConfig)

			if err != nil {
				panic(err)
			}

			err = dm.AddDevice(newDevice)

			if err != nil {
				panic(err)
			}
		} else {
			device.updateConfiguration(deviceConfig)
		}
	}

	for _, device := range dm.devices {
		dm.mutex.Lock()
		if _, found := midiConfigMap[device.GetAlias()]; !found {
			dm.mutex.Unlock()
			err := dm.RemoveDevice(device.GetAlias())
			if err != nil {
				continue
			}
		}
	}
}

func (dm *DeviceManager) GetSignals() chan core.Signal {
	return dm.signals
}

func NewDeviceManager() *DeviceManager {
	dm := DeviceManager{}
	dm.devices = make(map[string]*MidiDevice)
	dm.signals = make(chan core.Signal)

	return &dm
}
