package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"net"
	"os"
)

type RedisConfig struct {
	URL string `yaml:"url"`
}

type ManipulatorConfig struct {
	IPAddr string `yaml:"host"`
	Port   uint16 `yaml:"port"`
}

type ExecutorConfig struct {
	IPAddr string `yaml:"host"`
	Port   uint16 `yaml:"port"`
}

type MIDIConfig struct {
	DeviceName string  `yaml:"device_name"`
	HoldDelta  float64 `yaml:"hold_delta"`
}

type Config struct {
	RedisConfig       RedisConfig       `yaml:"redis"`
	ManipulatorConfig ManipulatorConfig `yaml:"manipulator"`
	ExecutorConfig    ExecutorConfig    `yaml:"executor"`
	MIDIConfig        MIDIConfig        `yaml:"midi"`
}

func (conf *Config) Validate() error {
	if conf.MIDIConfig.DeviceName == "" {
		return fmt.Errorf("valid MIDI device_name must be provided in config. Now {%s} is provided",
			conf.MIDIConfig.DeviceName)
	}
	if conf.MIDIConfig.HoldDelta < 0 {
		return fmt.Errorf("valid MIDI hold_delta must be provided in config. Now {%s} is provided",
			conf.MIDIConfig.DeviceName)
	}
	if conf.ManipulatorConfig.Port == 0 {
		return fmt.Errorf("valid manipulator agent port must be provided in config. Now {%d} is provided",
			conf.ManipulatorConfig.Port)
	}
	if manipulatorIP := net.ParseIP(conf.ManipulatorConfig.IPAddr); manipulatorIP == nil {
		return fmt.Errorf("valid manipulator agent ip must be provided in config. Now {%s} is provided",
			conf.ManipulatorConfig.IPAddr)
	}
	if conf.ExecutorConfig.Port == 0 {
		return fmt.Errorf("valid manipulator agent port must be provided in config. Now {%d} is provided",
			conf.ManipulatorConfig.Port)
	}
	if executorIP := net.ParseIP(conf.ExecutorConfig.IPAddr); executorIP == nil {
		return fmt.Errorf("valid manipulator agent ip must be provided in config. Now {%s} is provided",
			conf.ManipulatorConfig.IPAddr)
	}
	if conf.RedisConfig.URL == "" {
		return fmt.Errorf("valid Redis url must be provided in config. Now {%s} is provided",
			conf.RedisConfig.URL)
	}
	return nil
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

func ParseConfigFromBytes(data []byte) (*Config, error) {
	cfg := Config{}

	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
