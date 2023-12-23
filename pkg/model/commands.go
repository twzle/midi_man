package model

type MidiCommand interface {
	Code() string
	Description() string
}

type TurnLightOnCommand struct {
	KeyCode     int    `hubman:"key_code"`
	DeviceAlias string `hubman:"device_alias"`
}

func (c TurnLightOnCommand) Code() string {
	return "TurnLightOnCommand"
}

func (c TurnLightOnCommand) Description() string {
	return "Turns light on for specified MIDI key"
}

type TurnLightOffCommand struct {
	KeyCode     int    `hubman:"key_code"`
	DeviceAlias string `hubman:"device_alias"`
}

func (c TurnLightOffCommand) Code() string {
	return "TurnLightOffCommand"
}

func (c TurnLightOffCommand) Description() string {
	return "Turns light off for specified MIDI key"
}

type SingleBlinkCommand struct {
	KeyCode     int    `hubman:"key_code"`
	DeviceAlias string `hubman:"device_alias"`
}

func (c SingleBlinkCommand) Code() string {
	return "SingleBlinkCommand"
}

func (c SingleBlinkCommand) Description() string {
	return "Single blink (...->off->on->off) for specified MIDI key"
}

type SingleReversedBlinkCommand struct {
	KeyCode     int    `hubman:"key_code"`
	DeviceAlias string `hubman:"device_alias"`
}

func (c SingleReversedBlinkCommand) Code() string {
	return "SingleReversedBlinkCommand"
}

func (c SingleReversedBlinkCommand) Description() string {
	return "Single reverse blink (...->on->off->on) for specified MIDI key"
}

type ContinuousBlinkCommand struct {
	KeyCode     int    `hubman:"key_code"`
	DeviceAlias string `hubman:"device_alias"`
}

func (c ContinuousBlinkCommand) Code() string {
	return "ContinuousBlinkCommand"
}

func (c ContinuousBlinkCommand) Description() string {
	return "Continuous blink (until next discontinuous command) specified MIDI key"
}

type SetActiveNamespaceCommand struct {
	Namespace string `hubman:"namespace"`
	Device    string `hubman:"device"`
}

func (s SetActiveNamespaceCommand) Code() string {
	return "SetActiveNamespaceCommand"
}

func (s SetActiveNamespaceCommand) Description() string {
	return `Sets given namespace as active on given device, all signals will be received from will contain active namespace attribute`
}
