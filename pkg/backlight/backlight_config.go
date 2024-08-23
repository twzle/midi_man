package backlight

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Function initializing backlight configuration from YAML-compatible text file
func InitConfig(confPath string) (*DeviceBacklightConfig, error) {
	file, err := os.ReadFile(confPath) // O(N)
	if err != nil {
		return nil, err
	}
	cfg, err := ParseConfigFromBytes(file) // O(N)
	return cfg, err
}

// Function deserealizing backlight configuration
func ParseConfigFromBytes(data []byte) (*DeviceBacklightConfig, error) { // O(N)
	cfg := RawBacklightConfig{}

	err := yaml.Unmarshal(data, &cfg) // O(N)
	if err != nil {
		return nil, err
	}

	decodedCfg := decodeConfig(&cfg) // O(N)
	return &decodedCfg, nil
} // O(N)
