package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

// Default reconnect interval time
const MinReconnectIntervalMs = 1000

// Representation of configurtaion for trigger values used by set of controls
type TriggerValues struct {
	Increment int `json:"increment" yaml:"increment"`
	Decrement int `json:"decrement" yaml:"decrement"`
}

// Representation of configurtaion for single set of controls
type Controls struct {
	Keys         []int         `json:"keys" yaml:"keys"`
	Rotate       bool          `json:"rotate" yaml:"rotate"`
	ValueRange   [2]int        `json:"value_range" yaml:"value_range"`
	InitialValue int           `json:"initial_value" yaml:"initial_value"`
	Triggers     TriggerValues `json:"triggers" yaml:"triggers"`
}

// Representation of single device configurtaion
type DeviceConfig struct {
	DeviceName        string     `json:"device_name" yaml:"device_name"`
	StartupDelay      int        `json:"startup_delay" yaml:"startup_delay"`
	ReconnectInterval int        `json:"reconnect_interval" yaml:"reconnect_interval"`
	Active            bool       `json:"active" yaml:"active"`
	HoldDelta         int        `json:"hold_delta" yaml:"hold_delta"`
	Namespace         string     `json:"namespace" yaml:"namespace"`
	Controls          []Controls `json:"accumulate_controls" yaml:"accumulate_controls"`
	BlinkingPeriodMS  int        `json:"blinking_period_ms" yaml:"blinking_period_ms"`
}

// Representation of user configurtaion
type UserConfig struct {
	MidiDevices []DeviceConfig `json:"midi_devices" yaml:"midi_devices"`
}

// Function validating the contents of user configuration
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
				"Now {%s} is provided",
				idx, device.DeviceName, device.DeviceName)
		}
		if device.Namespace == "" {
			return fmt.Errorf("device #{%d} ({%s}) has no namespace specified", idx, device.DeviceName)
		}
		if device.HoldDelta < 0 {
			return fmt.Errorf(
				"device #{%d} ({%s}): hold_delta must be >=0ms. Now {%d} is provided",
				idx,
				device.DeviceName,
				device.HoldDelta,
			)
		}
		if device.StartupDelay < 0 {
			return fmt.Errorf(
				"device #{%d} ({%s}): startup_delay must be >=0ms. Now {%d} is provided",
				idx,
				device.DeviceName,
				device.StartupDelay,
			)
		}
		if device.ReconnectInterval < MinReconnectIntervalMs {
			return fmt.Errorf(
				"device #{%d} ({%s}): reconnect_interval must be >= %dms. Now {%d} is provided",
				idx,
				device.DeviceName,
				MinReconnectIntervalMs,
				device.ReconnectInterval,
			)
		}
		if device.BlinkingPeriodMS < 0 {
			return fmt.Errorf(
				"device #{%d} ({%s}): blinking_period_in_milliseconds must be >= 0ms. Now {%d} is provided",
				idx,
				device.DeviceName,
				device.ReconnectInterval,
			)
		}
	}
	return nil
}


// Function seraching duplicate device name in array of configured MIDI-devices
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

// Function deserealizing user configuration
func ParseConfigFromBytes(data []byte) (*UserConfig, error) {
	cfg := UserConfig{}

	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
