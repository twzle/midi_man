package backlight

import (
	"encoding/json"
	"os"
)

func InitConfig(confPath string) (*Raw_BacklightConfig, *Decoded_DeviceBacklightConfig, error) {
	jsonFile, err := os.ReadFile(confPath)
	if err != nil {
		return nil, nil, err
	}
	cfg, decodedCfg, err := ParseConfigFromBytes(jsonFile)
	return cfg, decodedCfg, err
}

func ParseConfigFromBytes(data []byte) (*Raw_BacklightConfig, *Decoded_DeviceBacklightConfig, error) {
	cfg := Raw_BacklightConfig{}

	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, nil, err
	}

	decodedCfg := DecodeConfig(&cfg)
	return &cfg, &decodedCfg, nil
}
