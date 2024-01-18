package midi

import (
	"os"
	"strings"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"go.uber.org/zap"
)

var healthCheckDelay = 400 * time.Millisecond

func CheckDevicesHealth(manager *DeviceManager) {
	for {
		manager.dMutex.Lock()
		for _, deviceName := range manager.deviceNames {
			if !HasDeviceWithName(deviceName, midi.GetInPorts()) {
				manager.logger.Error("InPort is unreachable for device", zap.String("alias", deviceName))
				os.Exit(1)
			}
			if !HasDeviceWithName(deviceName, midi.GetOutPorts()) {
				manager.logger.Error("OutPort is unreachable for device", zap.String("alias", deviceName))
				os.Exit(1)
			}
		}
		manager.dMutex.Unlock()
		time.Sleep(healthCheckDelay)
	}
}

func HasDeviceWithName[T drivers.Port](deviceName string, deviceList []T) bool {
	for _, portName := range deviceList {
		validPortName := GetValidPortName(portName.String())
		if strings.HasPrefix(validPortName, deviceName) {
			return true
		}
	}
	return false
}

func GetValidPortName(portName string) string {
	tokens := strings.Split(portName, " ")
	validPortName := strings.Join(tokens[:len(tokens)-1], " ")
	return validPortName
}
