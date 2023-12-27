package backlight

import "fmt"

func (db *DecodedDeviceBacklightConfig) TurnLight(deviceAlias string, key byte, color string, status string) ([]byte, error) {
	mapping, values := db.FindArguments(deviceAlias, key, color, status)

	if mapping == nil || values == nil {
		return nil, fmt.Errorf("parameters for TurnLight command were not found")
	}

	bytes := mapping.bytes

	/* Key takes single byte to be inserted into template byte sequence
	parsed from format string containing with key %key
	*/
	bytes[mapping.keyIdx] = key

	/* Payload takes multiple bytes to be inserted into template byte sequence
	parsed from format string containing with key %payload
	*/
	bytes = append(bytes[:mapping.payloadIdx+len(values.payload)-1], bytes[mapping.payloadIdx:]...)

	/* Insertion of payload byte sequence into template byte sequence starting from precalculated index
	 */
	for idx, b := range values.payload {
		bytes[mapping.payloadIdx+idx] = b
	}

	return bytes, nil
}
