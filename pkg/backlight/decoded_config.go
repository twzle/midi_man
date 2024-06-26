package backlight

import (
	"encoding/hex"
	"strings"
)

type StatusName string

const (
	On  StatusName = "on"
	Off StatusName = "off"
)

// Representation of decoded values
type DecodedValues struct {
	payload []byte
}

// Representation of decoded mapping
type DecodedMapping struct {
	payloadIdx     int
	keyIdx         int
	keyNumberShift int
	bytes          []byte
}

// Representation of decoded device backlight configuration
type DecodedDeviceBacklightConfig struct {
	ColorSetToValues          map[DecodedColorSetIdentifiers]DecodedValues
	KeyStatusToMapping        map[DecodedKeyStatusIdentifiers]DecodedMapping
	KeyBacklightMap           map[DecodedKeyBacklightIdentifiers]RawKeyBacklight
	DeviceKeyRangeMap         map[string][][2]byte
	DeviceBacklightTimeOffset map[string]int
}

// Representation of decoded color set identifiers
type DecodedColorSetIdentifiers struct {
	DeviceAlias string
	Status      StatusName
	ColorSpace  int
	ColorName   string
}

// Representation of decoded key status identifiers
type DecodedKeyStatusIdentifiers struct {
	DeviceAlias string
	Key         byte
	Status      StatusName
}

// Representation of decoded key backlight identifiers
type DecodedKeyBacklightIdentifiers struct {
	deviceAlias string
	key         byte
}

// Function decodes payload of given string
func decodePayload(payload string) []byte {
	byteString := strings.ReplaceAll(payload, " ", "")
	bytes, _ := hex.DecodeString(byteString)

	return bytes
}

// Function finds entry index by anchor in given string
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

// Function replaces anchors with values
func removeFormatKeysFromString(payload string) string {
	payload = strings.Replace(payload, "%payload", "00", 1)
	payload = strings.Replace(payload, "%key", "00", 1)
	return payload
}

// Function decodes mapping part of backlight configuration from raw format to optimized
func decodeMapping(byteString string, keyNumberShift int) DecodedMapping {
	payloadIdx := findEntryIndex(byteString, "%payload")
	keyIdx := findEntryIndex(byteString, "%key")
	payload := removeFormatKeysFromString(byteString)
	bytes := decodePayload(payload)
	return DecodedMapping{payloadIdx, keyIdx, keyNumberShift, bytes}

}

// Function decodes main part of backlight configuration from raw format to optimized
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
					On, deviceColorSpace.Id, onStatusColors.ColorName}

				values := DecodedValues{decodePayload(onStatusColors.Payload)}

				cstv[csi] = values
			}

			for _, offStatusColors := range deviceColorSpace.Off {

				csi := DecodedColorSetIdentifiers{deviceBacklightConfig.DeviceName,
					Off, deviceColorSpace.Id, offStatusColors.ColorName}

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
					key, On}

				kstm[ksi] = decodeMapping(backlightRange.BacklightStatuses.On.Bytes, backlightRange.KeyNumberShift)

				ksi = DecodedKeyStatusIdentifiers{deviceBacklightConfig.DeviceName,
					key, Off}

				kstm[ksi] = decodeMapping(backlightRange.BacklightStatuses.Off.Bytes, backlightRange.KeyNumberShift)

			}
		}
	}
	dbct := DecodedDeviceBacklightConfig{
		cstv, kstm, kbm,
		dkrm, dbto}
	return dbct
}

// Function finds arguments to deserealize backlight configuration
func (db *DecodedDeviceBacklightConfig) FindArguments(deviceAlias string, key byte, color string, status StatusName) (*DecodedMapping, *DecodedValues) {
	kbl := DecodedKeyBacklightIdentifiers{deviceAlias, key}
	kb, _ := db.KeyBacklightMap[kbl]

	csi := DecodedColorSetIdentifiers{
		deviceAlias, status, kb.ColorSpace, color}

	values, ok := db.ColorSetToValues[csi]

	if !ok {
		var fallbackColorName string

		switch status {
		case On:
			fallbackColorName = kb.BacklightStatuses.On.FallbackColor
		case Off:
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
