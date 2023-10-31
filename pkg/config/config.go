package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type RedisConfig struct {
	URL string `yaml:"url"`
}

type ManpulatorConfig struct {
	IPAddr string `yaml:"host"`
	Port   uint16 `yaml:"port"`
}

type MIDIConfig struct {
	DeviceName string  `yaml:"device_name"`
	HoldDelta  float32 `yaml:"hold_delta"`
}

type Config struct {
	RedisConfig      RedisConfig      `yaml:"redis"`
	ManpulatorConfig ManpulatorConfig `yaml:"manipulator"`
	MIDIConfig       MIDIConfig       `yaml:"midi"`
}

func InitConfig(confPath string) (*Config, error) {
	yFile, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	cfg := Config{}

	err = yaml.Unmarshal(yFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
