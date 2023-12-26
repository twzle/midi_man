package backlight

func TurnLightOn(config *Decoded_DeviceBacklightConfig, deviceAlias string, key int, color string) []byte {
	kbl := Decoded_KeyBacklightIdentifiers{deviceAlias, byte(key)}
	kb, _ := config.KeyBacklightMap[kbl]

	csi := Decoded_ColorSetIdentifiers{deviceAlias, "on", kb.ColorSpace, color}

	values, ok := config.ColorSetToValues[csi]

	if !ok {
		csi = Decoded_ColorSetIdentifiers{deviceAlias, "on", kb.ColorSpace,
			kb.BacklightStatuses.On.FallbackColor}
		values, ok = config.ColorSetToValues[csi]

		if !ok {
			return nil
		}
	}

	ksi := Decoded_KeyStatusIdentifiers{deviceAlias, byte(key), "on"}

	mapping := config.KeyStatusToMapping[ksi]

	bytes := mapping.bytes

	bytes[mapping.keyIdx] = byte(key)

	bytes = append(bytes[:mapping.payloadIdx+len(values.payload)-1], bytes[mapping.payloadIdx:]...)

	for idx, b := range values.payload {
		bytes[mapping.payloadIdx+idx] = b
	}

	return bytes
}

func TurnLightOff(config *Decoded_DeviceBacklightConfig, deviceAlias string, key int, color string) []byte {
	kbl := Decoded_KeyBacklightIdentifiers{deviceAlias, byte(key)}
	kb, _ := config.KeyBacklightMap[kbl]

	csi := Decoded_ColorSetIdentifiers{deviceAlias, "off", kb.ColorSpace, color}

	values, ok := config.ColorSetToValues[csi]

	if !ok {
		csi = Decoded_ColorSetIdentifiers{deviceAlias, "off", kb.ColorSpace,
			kb.BacklightStatuses.Off.FallbackColor}
		values, ok = config.ColorSetToValues[csi]

		if !ok {
			return nil
		}
	}

	ksi := Decoded_KeyStatusIdentifiers{deviceAlias, byte(key), "off"}

	mapping := config.KeyStatusToMapping[ksi]
	bytes := mapping.bytes

	bytes[mapping.keyIdx] = byte(key)

	bytes = append(bytes[:mapping.payloadIdx+len(values.payload)-1], bytes[mapping.payloadIdx:]...)
	for idx, b := range values.payload {
		bytes[mapping.payloadIdx+idx] = b
	}

	return bytes
}
