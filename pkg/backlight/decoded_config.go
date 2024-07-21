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
type Payload struct {
	payload []byte
}

// Representation of decoded mapping
type Mapping struct {
	payloadIdx     int
	keyIdx         int
	keyNumberShift int
	bytes          []byte
}

// Representation of decoded device backlight configuration
type DeviceBacklightConfig struct {
	ColorSetToValues          map[ColorSetIdentifiers]Payload
	KeyStatusToMapping        map[KeyStatusIdentifiers]Mapping
	KeyBacklightMap           map[KeyBacklightIdentifiers]RawKeyBacklight
	DeviceKeyRangeMap         map[string][][2]byte
	DeviceBacklightTimeOffset map[string]int
}

// Representation of decoded color set identifiers
type ColorSetIdentifiers struct {
	DeviceAlias string
	Status      StatusName
	ColorSpace  int
	ColorName   string
}

// Representation of decoded key status identifiers
type KeyStatusIdentifiers struct {
	DeviceAlias string
	Key         byte
	Status      StatusName
}

// Representation of decoded key backlight identifiers
type KeyBacklightIdentifiers struct {
	deviceAlias string
	key         byte
}

// Function decodes payload of given string
func decodePayload(payload string) []byte {
	byteString := strings.ReplaceAll(payload, " ", "") // O(N)
	bytes, _ := hex.DecodeString(byteString) // O(N)

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
	payload = strings.Replace(payload, "%payload", "00", 1) // O(N)
	payload = strings.Replace(payload, "%key", "00", 1) // O(N)
	return payload
}

// Function decodes mapping part of backlight configuration from raw format to optimized
func decodeMapping(byteString string, keyNumberShift int) Mapping { // O(N)
	payloadIdx := findEntryIndex(byteString, "%payload") // O(N)
	keyIdx := findEntryIndex(byteString, "%key") // O(N)
	payload := removeFormatKeysFromString(byteString) // O(N)
	bytes := decodePayload(payload) // O(N)
	return Mapping{payloadIdx, keyIdx, keyNumberShift, bytes}

}

// Function decodes main part of backlight configuration from raw format to optimized
func decodeConfig(cfg *RawBacklightConfig) DeviceBacklightConfig { // O(I + J) -> O(N)
	kbm := make(map[KeyBacklightIdentifiers]RawKeyBacklight)
	cstv := make(map[ColorSetIdentifiers]Payload)
	kstm := make(map[KeyStatusIdentifiers]Mapping)
	dkrm := make(map[string][][2]byte)
	dbto := make(map[string]int)

	for _, deviceBacklightConfig := range cfg.DeviceBacklightConfigurations { // N - устройств
		dbto[deviceBacklightConfig.DeviceName] = deviceBacklightConfig.BacklightTimeOffset

		for _, deviceColorSpace := range deviceBacklightConfig.ColorSpaces { // M - цветовых пространств
			for _, onStatusColors := range deviceColorSpace.On { // X - цветов включения

				csi := ColorSetIdentifiers{deviceBacklightConfig.DeviceName,
					On, deviceColorSpace.Id, onStatusColors.ColorName}

				values := Payload{decodePayload(onStatusColors.Payload)}

				cstv[csi] = values
			}

			for _, offStatusColors := range deviceColorSpace.Off { // Y - цветов выключения

				csi := ColorSetIdentifiers{deviceBacklightConfig.DeviceName,
					Off, deviceColorSpace.Id, offStatusColors.ColorName}

				values := Payload{decodePayload(offStatusColors.Payload)}

				cstv[csi] = values // ~O(1)
			}
		} // O(M) * O(max(X) + max(Y)) * O(1) = O(M * O(max(X) + max(Y)) * O(1)) = O(N)

		for _, backlightRange := range deviceBacklightConfig.KeyboardBacklight { // M - количество диапазонов клавиш
			keyRange := backlightRange.KeyRange

			dkrm[deviceBacklightConfig.DeviceName] = append(dkrm[deviceBacklightConfig.DeviceName], backlightRange.KeyRange)

			for key := keyRange[0]; key <= keyRange[len(keyRange)-1]; key++ { // O(2) - размер диапазона клавиш
				kbl := KeyBacklightIdentifiers{deviceBacklightConfig.DeviceName, key}
				kbm[kbl] = backlightRange

				ksi := KeyStatusIdentifiers{deviceBacklightConfig.DeviceName,
					key, On}

				kstm[ksi] = decodeMapping(backlightRange.BacklightStatuses.On.Bytes, backlightRange.KeyNumberShift) // O(X)

				ksi = KeyStatusIdentifiers{deviceBacklightConfig.DeviceName,
					key, Off}

				kstm[ksi] = decodeMapping(backlightRange.BacklightStatuses.Off.Bytes, backlightRange.KeyNumberShift) // O(Y)

			}
		} // O(M) * O(2) * O(max(X) + max(Y)) = O(M * 2 * max(X) + max(Y)) = O(N)
	} // O(I + J) -> O(N)
	dbct := DeviceBacklightConfig{
		cstv, kstm, kbm,
		dkrm, dbto}
	return dbct
}

// Function finds arguments to deserealize backlight configuration
func (db *DeviceBacklightConfig) FindArguments(deviceAlias string, key byte, color string, status StatusName) (*Mapping, *Payload) {
	kbl := KeyBacklightIdentifiers{deviceAlias, key}
	kb, _ := db.KeyBacklightMap[kbl] // O(1)

	csi := ColorSetIdentifiers{
		deviceAlias, status, kb.ColorSpace, color}

	values, ok := db.ColorSetToValues[csi] // O(1)

	if !ok {
		var fallbackColorName string

		switch status {
		case On:
			fallbackColorName = kb.BacklightStatuses.On.FallbackColor
		case Off:
			fallbackColorName = kb.BacklightStatuses.Off.FallbackColor
		}

		csi = ColorSetIdentifiers{
			deviceAlias, status, kb.ColorSpace, fallbackColorName}
		values, ok = db.ColorSetToValues[csi] // O(1)

		if !ok {
			return nil, nil
		}
	}

	ksi := KeyStatusIdentifiers{deviceAlias, key, status}

	mapping := db.KeyStatusToMapping[ksi] // O(1)

	return &mapping, &values
}
