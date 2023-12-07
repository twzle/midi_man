package core

import "midi_manipulator/pkg/utils"

func (dm *DeviceManager) TurnLightOnHandler(cmd utils.TurnLightOnCommand) {
	for _, device := range dm.devices {
		device.turnLightOn(cmd)
	}
}

func (dm *DeviceManager) TurnLightOffHandler(cmd utils.TurnLightOffCommand) {
	for _, device := range dm.devices {
		device.turnLightOff(cmd)
	}
}
