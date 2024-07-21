package midi

import (
	"fmt"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	_ "gitlab.com/gomidi/midi/v2"
	"go.uber.org/zap"
	"midi_manipulator/pkg/backlight"
	"midi_manipulator/pkg/config"
	"midi_manipulator/pkg/model"
	"sync"
)

// Representation of device manager entity
type DeviceManager struct {
	devices         map[string]*MidiDevice
	mutex           sync.Mutex
	signals         chan core.Signal
	backlightConfig *backlight.DeviceBacklightConfig
	logger          *zap.Logger
	checkManager    core.CheckRegistry
}

// Function sets current backlight configuration for devices
func (dm *DeviceManager) SetBacklightConfig(cfg *backlight.DeviceBacklightConfig) {
	dm.backlightConfig = cfg
}

// Function returns object representing MIDI-device from current device list
func (dm *DeviceManager) getDevice(alias string) (*MidiDevice, bool) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	device, ok := dm.devices[alias]
	return device, ok
}


// Function frees the resources of device manager with termination of all devices in current device list
func (dm *DeviceManager) Close() {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	for _, device := range dm.devices {
		device.Stop()
	}
}

// Function handles execution of command on device by its alias
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

// Function handles device list update process
func (dm *DeviceManager) UpdateDevices(midiConfig []config.DeviceConfig) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.logger.Info("Updating devices", zap.Int("deviceCount", len(midiConfig)), zap.Any("devices", midiConfig))
	dm.checkManager.Clear()

	for _, device := range dm.devices {
		device.Stop()
	}

	dm.devices = make(map[string]*MidiDevice)
	for _, deviceConfig := range midiConfig {
		newDevice := NewDevice(deviceConfig, dm.signals, dm.logger, dm.checkManager)
		dm.devices[newDevice.name] = newDevice
		go newDevice.RunDevice(dm.backlightConfig)
	}
}

// Function returns signal channel of device manager as value
func (dm *DeviceManager) GetSignals() chan core.Signal {
	return dm.signals
}

// Function initializing device manager entity with values
func NewDeviceManager(
	logger *zap.Logger,
	checkManager core.CheckRegistry,
) *DeviceManager {
	dm := DeviceManager{logger: logger, checkManager: checkManager}
	dm.devices = make(map[string]*MidiDevice)
	dm.signals = make(chan core.Signal)

	logger.Info("Created device manager")
	return &dm
}
