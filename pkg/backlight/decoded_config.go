package backlight

import (
	"encoding/hex"
	"strings"
)

type Decoded_Values struct {
	payload []byte
}

type Decoded_Mapping struct {
	payloadIdx int
	keyIdx     int
	bytes      []byte
}

type Decoded_DeviceBacklightConfig struct {
	ColorSetToValues   map[Decoded_ColorSetIdentifiers]Decoded_Values
	KeyStatusToMapping map[Decoded_KeyStatusIdentifiers]Decoded_Mapping
	KeyBacklightMap    map[Decoded_KeyBacklightIdentifiers]Raw_KeyBacklight
}

type Decoded_ColorSetIdentifiers struct {
	DeviceAlias string
	Status      string
	ColorSpace  int
	ColorName   string
}

type Decoded_KeyStatusIdentifiers struct {
	DeviceAlias string
	Key         byte
	Status      string
}

type Decoded_KeyBacklightIdentifiers struct {
	deviceAlias string
	key         byte
}

func DecodePayload(payload string) []byte {
	byteString := strings.ReplaceAll(payload, " ", "")
	bytes, _ := hex.DecodeString(byteString)

	return bytes
}

func FindEntryIndex(payload string, token string) int {
	idx := 0

	end := strings.Index(payload, token)
	for pos, char := range payload {
		if pos == end {
			break
		}

		if char == ' ' {
			idx++
		}
	}
	return idx
}

func RemoveFormatKeysFromString(payload string) string {
	payload = strings.Replace(payload, "%payload", "00", 1)
	payload = strings.Replace(payload, "%key", "00", 1)
	return payload
}

func DecodeMapping(byteString string) Decoded_Mapping {
	payloadIdx := FindEntryIndex(byteString, "%payload")
	keyIdx := FindEntryIndex(byteString, "%key")
	payload := RemoveFormatKeysFromString(byteString)
	bytes := DecodePayload(payload)
	return Decoded_Mapping{payloadIdx, keyIdx, bytes}

}

func DecodeConfig(cfg *Raw_BacklightConfig) Decoded_DeviceBacklightConfig {
	kbm := make(map[Decoded_KeyBacklightIdentifiers]Raw_KeyBacklight)
	cstv := make(map[Decoded_ColorSetIdentifiers]Decoded_Values)
	kstm := make(map[Decoded_KeyStatusIdentifiers]Decoded_Mapping)

	for _, deviceBacklightConfig := range cfg.DeviceBacklightConfigurations {
		for _, deviceColorSpace := range deviceBacklightConfig.ColorSpaces {
			for _, onStatusColors := range deviceColorSpace.On {

				csi := Decoded_ColorSetIdentifiers{deviceBacklightConfig.DeviceName,
					"on", deviceColorSpace.Id, onStatusColors.ColorName}

				values := Decoded_Values{DecodePayload(onStatusColors.Payload)}

				cstv[csi] = values
			}

			for _, offStatusColors := range deviceColorSpace.Off {

				csi := Decoded_ColorSetIdentifiers{deviceBacklightConfig.DeviceName,
					"off", deviceColorSpace.Id, offStatusColors.ColorName}

				values := Decoded_Values{DecodePayload(offStatusColors.Payload)}

				cstv[csi] = values
			}
		}

		for _, backlightRange := range deviceBacklightConfig.KeyboardBacklight {
			keyRange := backlightRange.KeyRange

			for key := keyRange[0]; key <= keyRange[len(keyRange)-1]; key++ {
				kbl := Decoded_KeyBacklightIdentifiers{deviceBacklightConfig.DeviceName, key}
				kbm[kbl] = backlightRange

				ksi := Decoded_KeyStatusIdentifiers{deviceBacklightConfig.DeviceName,
					key, "on"}

				kstm[ksi] = DecodeMapping(backlightRange.BacklightStatuses.On.Bytes)

				ksi = Decoded_KeyStatusIdentifiers{deviceBacklightConfig.DeviceName,
					key, "off"}

				kstm[ksi] = DecodeMapping(backlightRange.BacklightStatuses.Off.Bytes)

			}
		}
	}
	dbct := Decoded_DeviceBacklightConfig{cstv, kstm, kbm}
	return dbct
}
