package backlight

import (
	"encoding/json"
	"os"
)

func InitConfig(confPath string) (*RawBacklightConfig, *DecodedDeviceBacklightConfig, error) {
	jsonFile, err := os.ReadFile(confPath)
	if err != nil {
		return nil, nil, err
	}
	cfg, decodedCfg, err := ParseConfigFromBytes(jsonFile)
	return cfg, decodedCfg, err
}

func ParseConfigFromBytes(data []byte) (*RawBacklightConfig, *DecodedDeviceBacklightConfig, error) {
	cfg := RawBacklightConfig{}

	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, nil, err
	}

	decodedCfg := decodeConfig(&cfg)
	return &cfg, &decodedCfg, nil
}
