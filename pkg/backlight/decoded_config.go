package backlight

import (
	"encoding/hex"
	"strings"
)

type DecodedValues struct {
	payload []byte
}

type DecodedMapping struct {
	payloadIdx     int
	keyIdx         int
	keyNumberShift int
	bytes          []byte
}

type DecodedDeviceBacklightConfig struct {
	ColorSetToValues          map[DecodedColorSetIdentifiers]DecodedValues
	KeyStatusToMapping        map[DecodedKeyStatusIdentifiers]DecodedMapping
	KeyBacklightMap           map[DecodedKeyBacklightIdentifiers]RawKeyBacklight
	DeviceKeyRangeMap         map[string][][2]byte
	DeviceBacklightTimeOffset map[string]int
}

type DecodedColorSetIdentifiers struct {
	DeviceAlias string
	Status      string
	ColorSpace  int
	ColorName   string
}

type DecodedKeyStatusIdentifiers struct {
	DeviceAlias string
	Key         byte
	Status      string
}

type DecodedKeyBacklightIdentifiers struct {
	deviceAlias string
	key         byte
}

func decodePayload(payload string) []byte {
	byteString := strings.ReplaceAll(payload, " ", "")
	bytes, _ := hex.DecodeString(byteString)

	return bytes
}

func findEntryIndex(payload string, token string) int {
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

func removeFormatKeysFromString(payload string) string {
	payload = strings.Replace(payload, "%payload", "00", 1)
	payload = strings.Replace(payload, "%key", "00", 1)
	return payload
}

func decodeMapping(byteString string, keyNumberShift int) DecodedMapping {
	payloadIdx := findEntryIndex(byteString, "%payload")
	keyIdx := findEntryIndex(byteString, "%key")
	payload := removeFormatKeysFromString(byteString)
	bytes := decodePayload(payload)
	return DecodedMapping{payloadIdx, keyIdx, keyNumberShift, bytes}

}

func decodeConfig(cfg *RawBacklightConfig) DecodedDeviceBacklightConfig {
	kbm := make(map[DecodedKeyBacklightIdentifiers]RawKeyBacklight)
	cstv := make(map[DecodedColorSetIdentifiers]DecodedValues)
	kstm := make(map[DecodedKeyStatusIdentifiers]DecodedMapping)
	dkrm := make(map[string][][2]byte)
	dbto := make(map[string]int)

	for _, deviceBacklightConfig := range cfg.DeviceBacklightConfigurations {
		dbto[deviceBacklightConfig.DeviceName] = deviceBacklightConfig.BacklightTimeOffset

		for _, deviceColorSpace := range deviceBacklightConfig.ColorSpaces {
			for _, onStatusColors := range deviceColorSpace.On {

				csi := DecodedColorSetIdentifiers{deviceBacklightConfig.DeviceName,
					"on", deviceColorSpace.Id, onStatusColors.ColorName}

				values := DecodedValues{decodePayload(onStatusColors.Payload)}

				cstv[csi] = values
			}

			for _, offStatusColors := range deviceColorSpace.Off {

				csi := DecodedColorSetIdentifiers{deviceBacklightConfig.DeviceName,
					"off", deviceColorSpace.Id, offStatusColors.ColorName}

				values := DecodedValues{decodePayload(offStatusColors.Payload)}

				cstv[csi] = values
			}
		}

		for _, backlightRange := range deviceBacklightConfig.KeyboardBacklight {
			keyRange := backlightRange.KeyRange

			dkrm[deviceBacklightConfig.DeviceName] = append(dkrm[deviceBacklightConfig.DeviceName], backlightRange.KeyRange)

			for key := keyRange[0]; key <= keyRange[len(keyRange)-1]; key++ {
				kbl := DecodedKeyBacklightIdentifiers{deviceBacklightConfig.DeviceName, key}
				kbm[kbl] = backlightRange

				ksi := DecodedKeyStatusIdentifiers{deviceBacklightConfig.DeviceName,
					key, "on"}

				kstm[ksi] = decodeMapping(backlightRange.BacklightStatuses.On.Bytes, backlightRange.KeyNumberShift)

				ksi = DecodedKeyStatusIdentifiers{deviceBacklightConfig.DeviceName,
					key, "off"}

				kstm[ksi] = decodeMapping(backlightRange.BacklightStatuses.Off.Bytes, backlightRange.KeyNumberShift)

			}
		}
	}
	dbct := DecodedDeviceBacklightConfig{
		cstv, kstm, kbm,
		dkrm, dbto}
	return dbct
}

func (db *DecodedDeviceBacklightConfig) FindArguments(deviceAlias string, key byte, color string, status string) (*DecodedMapping, *DecodedValues) {
	kbl := DecodedKeyBacklightIdentifiers{deviceAlias, key}
	kb, _ := db.KeyBacklightMap[kbl]

	csi := DecodedColorSetIdentifiers{
		deviceAlias, status, kb.ColorSpace, color}

	values, ok := db.ColorSetToValues[csi]

	if !ok {
		var fallbackColorName string

		switch status {
		case "on":
			fallbackColorName = kb.BacklightStatuses.On.FallbackColor
		case "off":
			fallbackColorName = kb.BacklightStatuses.Off.FallbackColor
		}

		csi = DecodedColorSetIdentifiers{
			deviceAlias, status, kb.ColorSpace, fallbackColorName}
		values, ok = db.ColorSetToValues[csi]

		if !ok {
			return nil, nil
		}
	}

	ksi := DecodedKeyStatusIdentifiers{deviceAlias, key, status}

	mapping := db.KeyStatusToMapping[ksi]

	return &mapping, &values
}
