package pkg

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func main() {
	config, err := InitConfig("backlight.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	TurnLightOn(config, "MPD226", 4, "blue")

	TurnLightOff(config, "MPD226", 9, "light_red")
}

func KeyToByteString(key int) string {
	var b []byte
	b = append(b, byte(key))
	return hex.EncodeToString(b)
}

func GetSysexMsg(templateByteString string, key int, payload string) []byte {
	var sysexMsg []byte

	byteString := strings.ReplaceAll(templateByteString, "%key", KeyToByteString(key))
	byteString = strings.ReplaceAll(byteString, "%payload", payload)
	byteString = strings.ReplaceAll(byteString, " ", "")
	bytes, err := hex.DecodeString(byteString)

	if err != nil {
		fmt.Println(err)
	}

	sysexMsg = append(sysexMsg, bytes...)

	return sysexMsg
}

func GetDeviceBacklightConfig(config *BacklightConfig, deviceAlias string) *DeviceBacklightConfig {
	for _, deviceConfig := range config.DeviceBacklightConfigurations {
		if deviceConfig.DeviceName == deviceAlias {
			return &deviceConfig
		}
	}
	return nil
}

func GetKeyBacklight(config *DeviceBacklightConfig, key int) *KeyBacklight {
	for _, keyBacklightConfig := range config.KeyboardBacklight {
		if len(keyBacklightConfig.KeyRange) == 2 {
			if key >= keyBacklightConfig.KeyRange[0] && key <= keyBacklightConfig.KeyRange[1] {
				return &keyBacklightConfig
			}
		} else if len(keyBacklightConfig.KeyRange) == 1 {
			if keyBacklightConfig.KeyRange[0] == key {
				return &keyBacklightConfig
			}
		}
	}
	return nil
}

func GetColor(colorSpace *ColorSpace, colorName string, fallbackColorName string, status string) string {
	var colorSpaceStatus []Color
	if status == "On" {
		colorSpaceStatus = colorSpace.On
	} else {
		colorSpaceStatus = colorSpace.Off
	}

	for _, color := range colorSpaceStatus {
		if color.ColorName == colorName {
			return color.Payload
		}
	}

	for _, color := range colorSpaceStatus {
		if color.ColorName == fallbackColorName {
			return color.Payload
		}
	}

	return ""
}

func GetColorSpace(config *DeviceBacklightConfig, colorSpaceId int) *ColorSpace {
	for _, colorSpace := range config.ColorSpaces {
		if colorSpace.Id == colorSpaceId {
			return &colorSpace
		}
	}
	return nil
}

func GetMidiMessage(templateByteString string, key int, payload string) []byte {
	var cmd []byte
	cmd = GetSysexMsg(templateByteString, key, payload)

	return cmd
}

func TurnLightOn(config *BacklightConfig, deviceAlias string, key int, color string) {
	deviceBacklightConfig := GetDeviceBacklightConfig(config, deviceAlias)

	if deviceBacklightConfig == nil {
		return
	}

	keyBacklightConfig := GetKeyBacklight(deviceBacklightConfig, key)

	if keyBacklightConfig == nil {
		return
	}

	colorSpace := GetColorSpace(deviceBacklightConfig, keyBacklightConfig.ColorSpace)

	if colorSpace == nil {
		return
	}

	colorPayload := GetColor(colorSpace, color, keyBacklightConfig.BacklightStatuses.On.FallbackColor, "On")

	midiMsg := GetMidiMessage(keyBacklightConfig.BacklightStatuses.On.Bytes, key, colorPayload)

	fmt.Println(midiMsg)
}

func TurnLightOff(config *BacklightConfig, deviceAlias string, key int, color string) {
	deviceBacklightConfig := GetDeviceBacklightConfig(config, deviceAlias)

	if deviceBacklightConfig == nil {
		return
	}

	keyBacklightConfig := GetKeyBacklight(deviceBacklightConfig, key)

	if keyBacklightConfig == nil {
		return
	}

	colorSpace := GetColorSpace(deviceBacklightConfig, keyBacklightConfig.ColorSpace)

	if colorSpace == nil {
		return
	}

	colorPayload := GetColor(colorSpace, color, keyBacklightConfig.BacklightStatuses.Off.FallbackColor, "Off")

	midiMsg := GetMidiMessage(keyBacklightConfig.BacklightStatuses.Off.Bytes, key, colorPayload)

	fmt.Println(midiMsg)
}
