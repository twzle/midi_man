package backlight

import (
	"gopkg.in/yaml.v3"
	"os"
)

func InitConfig(confPath string) (*RawBacklightConfig, *DecodedDeviceBacklightConfig, error) {
	file, err := os.ReadFile(confPath)
	if err != nil {
		return nil, nil, err
	}
	cfg, decodedCfg, err := ParseConfigFromBytes(file)
	return cfg, decodedCfg, err
}

func ParseConfigFromBytes(data []byte) (*RawBacklightConfig, *DecodedDeviceBacklightConfig, error) {
	cfg := RawBacklightConfig{}

	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, nil, err
	}

	decodedCfg := decodeConfig(&cfg)
	return &cfg, &decodedCfg, nil
}
