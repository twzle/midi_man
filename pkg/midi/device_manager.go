package midi

import (
	"fmt"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	_ "gitlab.com/gomidi/midi/v2"
	"go.uber.org/zap"
	"log"
	"midi_manipulator/pkg/backlight"
	"midi_manipulator/pkg/config"
	"midi_manipulator/pkg/model"
	"sync"
)

type DeviceManager struct {
	devices         map[string]*MidiDevice
	mutex           sync.Mutex
	signals         chan core.Signal
	backlightConfig *backlight.DecodedDeviceBacklightConfig
	logger          *zap.Logger

	deviceNames []string
	dMutex      sync.Mutex
}

func (dm *DeviceManager) SetBacklightConfig(cfg *backlight.DecodedDeviceBacklightConfig) {
	dm.backlightConfig = cfg
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

func (dm *DeviceManager) addDeviceName(deviceName string) {
	dm.dMutex.Lock()
	defer dm.dMutex.Unlock()
	dm.deviceNames = append(dm.deviceNames, deviceName)
}

func (dm *DeviceManager) removeDevice(alias string) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	delete(dm.devices, alias)
}

func (dm *DeviceManager) removeDeviceName(alias string) {
	dm.dMutex.Lock()
	defer dm.dMutex.Unlock()
	for i, deviceName := range dm.deviceNames {
		if deviceName == alias {
			dm.deviceNames = append(dm.deviceNames[:i], dm.deviceNames[i+1:]...)
		}
	}
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
		dm.logger.Warn("Received command for non existing device", zap.String("device", alias))
		return fmt.Errorf("device with alias {%s} doesn't exist", alias)
	}

	if !device.active { // Possibly should be deprecated?
		dm.logger.Warn("Received command for inactive device", zap.String("device", alias))
		return fmt.Errorf("device with alias {%s} is not active", alias)
	}

	err := device.ExecuteCommand(cmd, dm.backlightConfig)

	if err != nil {
		return err
	}

	return nil
}

func (dm *DeviceManager) AddDevice(device *MidiDevice) error {
	dm.logger.Info("Adding device", zap.String("alias", device.GetAlias()))
	_, found := dm.getDevice(device.GetAlias())

	if found {
		return fmt.Errorf("device {%s} already exists", device.GetAlias())
	}

	dm.addDevice(device)
	dm.addDeviceName(device.GetAlias())
	err := device.RunDevice(dm.signals, dm.backlightConfig)

	if err != nil {
		return err
	}

	return nil
}

func (dm *DeviceManager) RemoveDevice(alias string) error {
	dm.logger.Info("Removing device", zap.String("device", alias))

	device, found := dm.getDevice(alias)
	if !found {
		return fmt.Errorf("device {%s} doesn't exist", device.GetAlias())
	}

	err := device.StopDevice()

	if err != nil {
		return err
	}

	dm.removeDevice(alias)
	dm.removeDeviceName(alias)

	return nil
}

func (dm *DeviceManager) UpdateDevices(midiConfig []config.DeviceConfig) {
	dm.logger.Info("Updating devices", zap.Int("deviceCount", len(midiConfig)))
	var midiConfigMap = make(map[string]config.DeviceConfig)

	for _, device := range midiConfig {
		midiConfigMap[device.DeviceName] = device
	}

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
			device.updateConfiguration(deviceConfig, dm.signals)
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

func (dm *DeviceManager) GetSignals() chan core.Signal {
	return dm.signals
}

func (dm *DeviceManager) SetActiveNamespace(newActive string, device string) {
	dm.logger.Info("Setting namespace as active", zap.String("namespace", newActive))
	dm.mutex.Lock()
	d, ok := dm.devices[device]
	if !ok {
		dm.logger.Error("Not found given device for namespace", zap.String("device", device), zap.String("newNamespace", newActive))
	} else {
		d.mutex.Lock()
		d.namespace = newActive
		d.mutex.Unlock()
		d.sendNamespaceChangedSignal(dm.signals)
	}
	dm.mutex.Unlock()
}

func NewDeviceManager(logger *zap.Logger) *DeviceManager {
	dm := DeviceManager{logger: logger}
	dm.devices = make(map[string]*MidiDevice)
	dm.signals = make(chan core.Signal)

	logger.Info("Created device manager")
	return &dm
}
