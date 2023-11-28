package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type RedisConfig struct {
	URL string `json:"url"`
}

type AppConfig struct {
	IPAddr string `json:"host"`
	Port   uint16 `json:"port"`
}

type ExecutorConfig struct {
	IPAddr string `json:"host"`
	Port   uint16 `json:"port"`
}

type MidiConfig struct {
	DeviceName string  `json:"device_name"`
	HoldDelta  float64 `json:"hold_delta"`
}

type Config struct {
	RedisConfig RedisConfig  `json:"redis"`
	AppConfig   AppConfig    `json:"app"`
	MidiConfig  []MidiConfig `json:"midi_devices"`
}

func (conf *Config) Validate() error {
	if len(conf.MidiConfig) == 0 {
		return fmt.Errorf("MIDI devices were not found in configuration file")
	}
	for idx, device := range conf.MidiConfig {
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
	if conf.AppConfig.Port == 0 {
		return fmt.Errorf("valid manipulator agent port must be provided in config. Now {%d} is provided",
			conf.AppConfig.Port)
	}
	if manipulatorIP := net.ParseIP(conf.AppConfig.IPAddr); manipulatorIP == nil {
		return fmt.Errorf("valid manipulator agent ip must be provided in config. Now {%s} is provided",
			conf.AppConfig.IPAddr)
	}
	if conf.RedisConfig.URL == "" {
		return fmt.Errorf("valid Redis url must be provided in config. Now {%s} is provided",
			conf.RedisConfig.URL)
	}
	return nil
}

func InitConfig(confPath string) (*Config, error) {
	jsonFile, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	cfg, err := ParseConfigFromBytes(jsonFile)
	return cfg, err
}

func ParseConfigFromBytes(data []byte) (*Config, error) {
	cfg := Config{}

	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
