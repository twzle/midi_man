package midi

import (
	"gitlab.com/gomidi/midi/v2/drivers"
	"strings"
)

func HasDeviceWithName[T drivers.Port](deviceName string, deviceList []T) bool {
	for _, portName := range deviceList {
		if strings.HasPrefix(portName.String(), deviceName) {
			return true
		}
	}
	return false
}
