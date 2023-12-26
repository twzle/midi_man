package backlight

func On(key byte, alias string, color string) {
	dbct := Decoded_DeviceBacklightConfigTree{}
	kb, _ := dbct.keyBacklightMap[key]

	csi := Decoded_ColorSetIdentifiers{alias, "on", kb.ColorSpace, color}

	values, ok := dbct.colorSetToValues[csi] // {"payload": "B1 00 7F", "key": 9}

	if !ok {
		csi = Decoded_ColorSetIdentifiers{alias, "on", kb.ColorSpace,
			kb.BacklightStatuses.On.FallbackColor}
		values, ok = dbct.colorSetToValues[csi]

		if !ok {
			return
		}
	}

	ksi := Decoded_KeyStatusIdentifiers{alias, key, "on"}

	mapping := dbct.keyStatusToMapping[ksi]
	bytes := mapping.bytes

	bytes = append(bytes[:mapping.payloadIdx+len(values.payload)], bytes[mapping.payloadIdx:]...)

	for idx, b := range mapping.bytes {
		bytes[mapping.payloadIdx+idx] = b
	}

	bytes = append(bytes[:mapping.keyIdx+1], bytes[mapping.keyIdx:]...)
	bytes[mapping.keyIdx] = key
}
