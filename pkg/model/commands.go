package model

// Represantion of common interface for commands
type MidiCommand interface {
	Code() string
	Description() string
}

// Representation of command to turn on the backlight of single component
type TurnLightOnCommand struct {
	KeyCode     int    `hubman:"key_code"`
	DeviceAlias string `hubman:"device_alias"`
	ColorName   string `hubman:"color_name"`
}

// Function returns string representation of model
func (c TurnLightOnCommand) Code() string {
	return "TurnLightOnCommand"
}

// Function returns string description of model
func (c TurnLightOnCommand) Description() string {
	return "Turns light on for specified MIDI key"
}

// Representation of command to turn off the backlight of single component
type TurnLightOffCommand struct {
	KeyCode     int    `hubman:"key_code"`
	DeviceAlias string `hubman:"device_alias"`
	ColorName   string `hubman:"color_name"`
}

// Function returns string representation of model
func (c TurnLightOffCommand) Code() string {
	return "TurnLightOffCommand"
}

// Function returns string description of model
func (c TurnLightOffCommand) Description() string {
	return "Turns light off for specified MIDI key"
}

// Representation of command to call backlight blink of single component
type SingleBlinkCommand struct {
	KeyCode     int    `hubman:"key_code"`
	DeviceAlias string `hubman:"device_alias"`
	ColorName   string `hubman:"color_name"`
}

// Function returns string representation of model
func (c SingleBlinkCommand) Code() string {
	return "SingleBlinkCommand"
}

// Function returns string description of model
func (c SingleBlinkCommand) Description() string {
	return "Single blink (...->off->on->off) for specified MIDI key"
}

// Representation of command to call backlight reversed blink of single component
type SingleReversedBlinkCommand struct {
	KeyCode     int    `hubman:"key_code"`
	DeviceAlias string `hubman:"device_alias"`
	ColorName   string `hubman:"color_name"`
}

// Function returns string representation of model
func (c SingleReversedBlinkCommand) Code() string {
	return "SingleReversedBlinkCommand"
}

// Function returns string description of model
func (c SingleReversedBlinkCommand) Description() string {
	return "Single reverse blink (...->on->off->on) for specified MIDI key"
}

// Representation of command to call continuous backlight blink of single component
type ContinuousBlinkCommand struct {
	KeyCode     int    `hubman:"key_code"`
	DeviceAlias string `hubman:"device_alias"`
	ColorName   string `hubman:"color_name"`
}

// Function returns string representation of model
func (c ContinuousBlinkCommand) Code() string {
	return "ContinuousBlinkCommand"
}

// Function returns string description of model
func (c ContinuousBlinkCommand) Description() string {
	return "Continuous blink (until next discontinuous command) specified MIDI key"
}

// Representation of command to set active namespace of single device
type SetActiveNamespaceCommand struct {
	Namespace   string `hubman:"namespace"`
	DeviceAlias string `hubman:"device"`
}

// Function returns string representation of model
func (s SetActiveNamespaceCommand) Code() string {
	return "SetActiveNamespaceCommand"
}

// Function returns string description of model
func (s SetActiveNamespaceCommand) Description() string {
	return `Sets given namespace as active on given device, all signals will be received from will contain active namespace attribute`
}

// Representation of command to start backlight blinking of single component
type StartBlinkingCommand struct {
	KeyCode      int    `hubman:"key_code"`
	DeviceAlias  string `hubman:"device_alias"`
	OnColorName  string `hubman:"on_color_name"`
	OffColorName string `hubman:"off_color_name"`
}

// Function returns string representation of model
func (s StartBlinkingCommand) Code() string {
	return "StartBlinkingCommand"
}

// Function returns string description of model
func (s StartBlinkingCommand) Description() string {
	return "Make the key repeatedly blink with chosen color"
}

// Representation of command to stop backlight blinking of single component
type StopBlinkingCommand struct {
	KeyCode     int    `hubman:"key_code"`
	DeviceAlias string `hubman:"device_alias"`
}

// Function returns string representation of model
func (s StopBlinkingCommand) Code() string {
	return "StopBlinkingCommand"
}

// Function returns string description of model
func (s StopBlinkingCommand) Description() string {
	return "Make the key stop blinking if it blinks"
}
