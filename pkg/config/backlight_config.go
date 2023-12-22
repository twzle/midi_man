package pkg

import (
	"encoding/json"
	"os"
)

type Color struct {
	ColorName string `json:"color_name" yaml:"color_name"`
	Payload   string `json:"payload" yaml:"payload"`
}

type ColorSpace struct {
	Id  int     `json:"color_space_id" yaml:"color_space_id"`
	On  []Color `json:"on" yaml:"on"`
	Off []Color `json:"off" yaml:"off"`
}

type Status struct {
	Type          string `json:"type" yaml:"type"`
	FallbackColor string `json:"fallback_color" yaml:"fallback_color"`
	Bytes         string `json:"bytes" yaml:"bytes"`
}

type KeyBacklightStatuses struct {
	On  Status `json:"on" yaml:"on"`
	Off Status `json:"off" yaml:"off"`
}

type KeyBacklight struct {
	KeyRange          []int                `json:"key_range" yaml:"key_range"`
	ColorSpace        int                  `json:"color_space" yaml:"color_space"`
	BacklightStatuses KeyBacklightStatuses `json:"statuses" yaml:"statuses"`
}

type DeviceBacklightConfig struct {
	DeviceName        string         `json:"device_name" yaml:"device_name"`
	ColorSpaces       []ColorSpace   `json:"color_spaces" yaml:"color_spaces"`
	KeyboardBacklight []KeyBacklight `json:"keyboard_backlight" yaml:"keyboard_backlight"`
}

type BacklightConfig struct {
	DeviceBacklightConfigurations []DeviceBacklightConfig `json:"device_light_configuration" yaml:"device_light_configuration"`
}

func InitConfig(confPath string) (*BacklightConfig, error) {
	jsonFile, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	cfg, err := ParseConfigFromBytes(jsonFile)
	return cfg, err
}

func ParseConfigFromBytes(data []byte) (*BacklightConfig, error) {
	cfg := BacklightConfig{}

	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
