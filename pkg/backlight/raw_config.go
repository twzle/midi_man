package main

import (
	"encoding/json"
	"os"
)

type Raw_Color struct {
	ColorName string `json:"color_name" yaml:"color_name"`
	Payload   string `json:"payload" yaml:"payload"`
}

type Raw_ColorSpace struct {
	Id  int         `json:"color_space_id" yaml:"color_space_id"`
	On  []Raw_Color `json:"on" yaml:"on"`
	Off []Raw_Color `json:"off" yaml:"off"`
}

type Raw_Status struct {
	Type          string `json:"type" yaml:"type"`
	FallbackColor string `json:"fallback_color" yaml:"fallback_color"`
	Bytes         string `json:"bytes" yaml:"bytes"`
}

type Raw_KeyBacklightStatuses struct {
	On  Raw_Status `json:"on" yaml:"on"`
	Off Raw_Status `json:"off" yaml:"off"`
}

type Raw_KeyBacklight struct {
	KeyRange          []byte                   `json:"key_range" yaml:"key_range"`
	ColorSpace        int                      `json:"color_space" yaml:"color_space"`
	BacklightStatuses Raw_KeyBacklightStatuses `json:"statuses" yaml:"statuses"`
}

type Raw_DeviceBacklightConfig struct {
	DeviceName        string             `json:"device_name" yaml:"device_name"`
	ColorSpaces       []Raw_ColorSpace   `json:"color_spaces" yaml:"color_spaces"`
	KeyboardBacklight []Raw_KeyBacklight `json:"keyboard_backlight" yaml:"keyboard_backlight"`
}

type Raw_BacklightConfig struct {
	DeviceBacklightConfigurations []Raw_DeviceBacklightConfig `json:"device_light_configuration" yaml:"device_light_configuration"`
}

func InitConfig(confPath string) (*Raw_BacklightConfig, error) {
	jsonFile, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	cfg, err := ParseConfigFromBytes(jsonFile)
	return cfg, err
}

func ParseConfigFromBytes(data []byte) (*Raw_BacklightConfig, error) {
	cfg := Raw_BacklightConfig{}

	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	DecodeConfig(&cfg)
	return &cfg, nil
}
