package backlight

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Function initializing backlight configuration from YAML-compatible text file
func InitConfig(confPath string) (*DecodedDeviceBacklightConfig, error) {
	file, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	cfg, err := ParseConfigFromBytes(file)
	return cfg, err
}

// Function deserealizing backlight configuration
func ParseConfigFromBytes(data []byte) (*DecodedDeviceBacklightConfig, error) {
	cfg := RawBacklightConfig{}

	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	decodedCfg := decodeConfig(&cfg)
	return &decodedCfg, nil
}
