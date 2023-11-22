package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type OnStatus struct {
	MidiType     string `json:"type"`
	MidiChannel  int    `json:"channel"`
	MidiVelocity int    `json:"velocity"`
}

type OffStatus struct {
	MidiType     string `json:"type"`
	MidiChannel  int    `json:"channel"`
	MidiVelocity int    `json:"velocity"`
}

type Statuses struct {
	On  OnStatus  `json:"on"`
	Off OffStatus `json:"off"`
}

type Backlight struct {
	Key               int      `json:"key"`
	IsRGB             bool     `json:"is_rgb"`
	HasMultipleColors bool     `json:"has_multiple_colors"`
	Statuses          Statuses `json:"statuses"`
}

type BacklightConfig struct {
	Backlight []Backlight `json:"backlight"`
}

func InitBacklightConfig(deviceName string) (*BacklightConfig, error) {
	confPath := fmt.Sprintf("configs/backlight/%s.json", deviceName)
	jsonFile, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	cfg, err := ParseBacklightConfigFromBytes(jsonFile)
	return cfg, err
}

func ParseBacklightConfigFromBytes(data []byte) (*BacklightConfig, error) {
	cfg := BacklightConfig{}

	err := json.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &cfg, nil
}
