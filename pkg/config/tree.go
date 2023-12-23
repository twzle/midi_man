package main

import (
	"github.com/rdleal/intervalst/interval"
)

type Values struct {
	payload []byte
	key     int
}

type Mapping struct {
	payloadIdx int
	keyIdx     int
	bytes      []byte
}

type DeviceBacklightConfigTree struct {
	colorSetToValues   map[ColorSetIdentifiers]Values
	keyStatusToMapping map[KeyStatusIdentifiers]Mapping
	keyRangeTree       interval.SearchTree[KeyBacklight, int]
}

type ColorSetIdentifiers struct {
	DeviceAlias string
	Status      string
	ColorSpace  int
	ColorName   string
}

type KeyStatusIdentifiers struct {
	DeviceAlias string
	Key         int
	Status      string
}

// Get slices of bytes for midi message on received TurnOnLightCommand
func on() {
	var key = 9
	var alias = "MPD226"
	var color = "red"

	dbct := DeviceBacklightConfigTree{}
	kb, _ := dbct.keyRangeTree.Find(key, key)

	csi := ColorSetIdentifiers{alias, "on", kb.ColorSpace, color}

	values, ok := dbct.colorSetToValues[csi] // {"payload": "B1 00 7F", "key": 9}

	if !ok {
		csi = ColorSetIdentifiers{alias, "on", kb.ColorSpace,
			kb.BacklightStatuses.On.FallbackColor}
		values, ok = dbct.colorSetToValues[csi]

		if !ok {
			return
		}
	}

	ksi := KeyStatusIdentifiers{alias, key, "on"}

	mapping := dbct.keyStatusToMapping[ksi]
	bytes := mapping.bytes

	bytes = append(bytes[:mapping.payloadIdx+len(values.payload)], bytes[mapping.payloadIdx:]...)

	for idx, b := range mapping.bytes {
		bytes[mapping.payloadIdx+idx] = b
	}

	bytes = append(bytes[:mapping.keyIdx+1], bytes[mapping.keyIdx:]...)
	bytes[mapping.keyIdx] = byte(values.key)
}
