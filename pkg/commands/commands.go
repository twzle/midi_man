package commands

type TurnLightOnCommand struct {
	KeyCode     float64 `hubman:"key_code"`
	RuleCommand string  `hubman:"rule_command"`
}

func (c TurnLightOnCommand) Code() string {
	return "TurnLightOnCommand"
}

func (c TurnLightOnCommand) Description() string {
	return "Turns light on for specified MIDI key"
}

type TurnLightOffCommand struct {
	KeyCode     float64 `hubman:"key_code"`
	RuleCommand string  `hubman:"rule_command"`
}

func (c TurnLightOffCommand) Code() string {
	return "TurnLightOffCommand"
}

func (c TurnLightOffCommand) Description() string {
	return "Turns light off for specified MIDI key"
}

type SingleBlinkCommand struct {
	KeyCode     float64 `hubman:"key_code"`
	RuleCommand string  `hubman:"rule_command"`
}

func (c SingleBlinkCommand) Code() string {
	return "SingleBlinkCommand"
}

func (c SingleBlinkCommand) Description() string {
	return "Single blink (...->off->on->off) for specified MIDI key"
}

type SingleReversedBlinkCommand struct {
	KeyCode     float64 `hubman:"key_code"`
	RuleCommand string  `hubman:"rule_command"`
}

func (c SingleReversedBlinkCommand) Code() string {
	return "SingleReversedBlinkCommand"
}

func (c SingleReversedBlinkCommand) Description() string {
	return "Single reverse blink (...->on->off->on) for specified MIDI key"
}

type ContinuousBlinkCommand struct {
	KeyCode     float64 `hubman:"key_code"`
	RuleCommand string  `hubman:"rule_command"`
}

func (c ContinuousBlinkCommand) Code() string {
	return "ContinuousBlinkCommand"
}

func (c ContinuousBlinkCommand) Description() string {
	return "Continuous blink (until next discontinuous command) specified MIDI key"
}
