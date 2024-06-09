package midi

import (
	"gitlab.com/gomidi/midi/v2/drivers"
	"strings"
)

// Function checks if device is present in device list
func HasDeviceWithName[T drivers.Port](deviceName string, deviceList []T) bool {
	for _, portName := range deviceList {
		if strings.Contains(portName.String(), deviceName) || strings.Contains(deviceName, portName.String()) {
			return true
		}
	}
	return false
}
