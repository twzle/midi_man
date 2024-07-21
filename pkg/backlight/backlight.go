package backlight

import (
	"errors"
)

// Function contains logic to perform backlight operations with MIDI-devices using backlight configuration
func (db *DeviceBacklightConfig) TurnLight(deviceAlias string, key byte, color string, status StatusName) ([]byte, error) { // O(N)
	mapping, values := db.FindArguments(deviceAlias, key, color, status) // O(1)

	if mapping == nil || values == nil {
		// TODO: create global variable for errors.New in order to validate outside
		return nil, errors.New("parameters for TurnLight command were not found")
	}

	bytes := make([]byte, len(mapping.bytes), cap(mapping.bytes)) // O(N)
	copy(bytes, mapping.bytes) // O(N)

	/* Key takes single byte to be inserted into template byte sequence
	parsed from format string containing with key %key
	*/
	bytes[mapping.keyIdx] = byte(int(key) - mapping.keyNumberShift) // O(1)

	/* Payload takes multiple bytes to be inserted into template byte sequence
	parsed from format string containing with key %payload
	*/
	bytes = append(bytes[:mapping.payloadIdx+len(values.payload)-1], bytes[mapping.payloadIdx:]...) // O(N)

	/* Insertion of payload byte sequence into template byte sequence starting from precalculated index
	 */
	for idx, b := range values.payload { // O(N)
		bytes[mapping.payloadIdx+idx] = b
	}

	return bytes, nil
} // O(N)
