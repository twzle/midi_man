package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type TriggerValues struct {
	Increment byte `json:"increment" yaml:"increment"`
	Decrement byte `json:"decrement" yaml:"decrement"`
}

type Controls struct {
	Keys         []byte        `json:"keys" yaml:"keys"`
	Rotate       bool          `json:"rotate" yaml:"rotate"`
	ValueRange   [2]byte       `json:"value_range" yaml:"value_range"`
	InitialValue byte          `json:"initial_value" yaml:"initial_value"`
	Triggers     TriggerValues `json:"triggers" yaml:"triggers"`
}

type DeviceConfig struct {
	DeviceName string   `json:"device_name" yaml:"device_name"`
	Active     bool     `json:"active" yaml:"active"`
	HoldDelta  float64  `json:"hold_delta" yaml:"hold_delta"`
	Namespace  string   `json:"namespace" yaml:"namespace"`
	Controls   Controls `json:"accumulate_controls" yaml:"accumulate_controls"`
}

type UserConfig struct {
	MidiDevices []DeviceConfig `json:"midi_devices" yaml:"midi_devices"`
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
		if device.Namespace == "" {
			return fmt.Errorf("device #{%d} ({%s}) has no namespace specified", idx, device.DeviceName)
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

func ParseConfigFromBytes(data []byte) (*UserConfig, error) {
	cfg := UserConfig{}

	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
