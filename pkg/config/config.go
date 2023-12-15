package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type DeviceConfig struct {
	DeviceName string  `yaml:"device_name"`
	Active     bool    `yaml:"active"`
	HoldDelta  float64 `yaml:"hold_delta"`
}

type UserConfig struct {
	MidiDevices []DeviceConfig `yaml:"midi_devices"`
}

func (conf *UserConfig) Validate() error {
	if len(conf.MidiDevices) == 0 {
		fmt.Println("MIDI devices were not found in configuration file")
	}
	if alias, has := conf.hasDuplicateDevices(); has {
		return fmt.Errorf("found duplicate MIDI device with alias {%s} in config", alias)
	}
	for idx, device := range conf.MidiDevices {
		if device.DeviceName == "" {
			return fmt.Errorf("device #{%d} ({%s}): "+
				"valid MIDI device_name must be provided in config. "+
				"Now {%f} is provided",
				idx, device.DeviceName, device.HoldDelta)
		}
		if device.HoldDelta < 0 {
			return fmt.Errorf("device #{%d} ({%s}): "+
				"valid MIDI hold_delta must be provided in config."+
				" Now {%f} is provided",
				idx, device.DeviceName, device.HoldDelta)
		}
	}
	return nil
}

func (conf *UserConfig) hasDuplicateDevices() (string, bool) {
	x := make(map[string]struct{})

	for _, v := range conf.MidiDevices {
		if _, has := x[v.DeviceName]; has {
			return v.DeviceName, true
		}
		x[v.DeviceName] = struct{}{}
	}

	return "", false
}

func InitConfig(confPath string) (*UserConfig, error) {
	jsonFile, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	cfg, err := ParseConfigFromBytes(jsonFile)
	return cfg, err
}

func ParseConfigFromBytes(data []byte) (*UserConfig, error) {
	cfg := UserConfig{}

	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
