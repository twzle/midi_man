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

type DeviceManager struct {
	devices         map[string]*MidiDevice
	mutex           sync.Mutex
	signals         chan core.Signal
	backlightConfig *backlight.DecodedDeviceBacklightConfig
	logger          *zap.Logger
	checkManager    core.CheckRegistry
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

func (dm *DeviceManager) Close() {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	for _, device := range dm.devices {
		device.Stop()
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
		newDevice.RunDevice(dm.backlightConfig)
	}
}

func (dm *DeviceManager) GetSignals() chan core.Signal {
	return dm.signals
}

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
