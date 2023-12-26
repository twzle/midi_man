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

type Decoded_DeviceBacklightConfigTree struct {
	colorSetToValues   map[Decoded_ColorSetIdentifiers]Decoded_Values
	keyStatusToMapping map[Decoded_KeyStatusIdentifiers]Decoded_Mapping
	keyBacklightMap    map[byte]Raw_KeyBacklight
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
	payload = strings.Replace(payload, "%payload", "", 1)
	payload = strings.Replace(payload, "%key", "", 1)
	return payload
}

func DecodeMapping(byteString string) Decoded_Mapping {
	payloadIdx := FindEntryIndex(byteString, "%payload")
	keyIdx := FindEntryIndex(byteString, "%key")
	payload := RemoveFormatKeysFromString(byteString)
	bytes := DecodePayload(payload)
	return Decoded_Mapping{payloadIdx, keyIdx, bytes}

}

func DecodeConfig(cfg *Raw_BacklightConfig) Decoded_DeviceBacklightConfigTree {
	kbm := make(map[byte]Raw_KeyBacklight)
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
			for _, key := range backlightRange.KeyRange {
				kbm[key] = backlightRange

				ksi := Decoded_KeyStatusIdentifiers{deviceBacklightConfig.DeviceName,
					key, "on"}

				kstm[ksi] = DecodeMapping(backlightRange.BacklightStatuses.On.Bytes)

				ksi = Decoded_KeyStatusIdentifiers{deviceBacklightConfig.DeviceName,
					key, "off"}

				kstm[ksi] = DecodeMapping(backlightRange.BacklightStatuses.Off.Bytes)

			}
		}
	}
	dbct := Decoded_DeviceBacklightConfigTree{cstv, kstm, kbm}
	return dbct
}
