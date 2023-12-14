package model

type NotePushed struct {
	Device   string `hubman:"device"`
	KeyCode  int    `hubman:"key_code"`
	Velocity int    `hubman:"velocity"`
}

func (s NotePushed) Code() string {
	return "NotePushed"
}

func (s NotePushed) Description() string {
	return "NotePushed - signal represents state of key with 'Note' type right off it was pressed on a device"
}

type NoteHold struct {
	Device   string `hubman:"device"`
	KeyCode  int    `hubman:"key_code"`
	Velocity int    `hubman:"velocity"`
}

func (s NoteHold) Code() string {
	return "NoteHold"
}

func (s NoteHold) Description() string {
	return "NoteHold - signal represents state of key with 'Note' type that is pressed for long"
}

type NoteReleased struct {
	Device   string `hubman:"device"`
	KeyCode  int    `hubman:"key_code"`
	Velocity int    `hubman:"velocity"`
}

func (s NoteReleased) Code() string {
	return "NoteReleased"
}

func (s NoteReleased) Description() string {
	return "NoteReleased - signal represents state of key with 'Note' type right off it was released on a device"
}

// ControlPushed В MIDI Control имеет только один тип событий "ControlChange",
// поэтому длительность и конец нажатия здесь не отслеживаются
type ControlPushed struct {
	Device  string `hubman:"device"`
	KeyCode int    `hubman:"key_code"`
	Value   int    `hubman:"velocity"`
}

func (s ControlPushed) Code() string {
	return "ControlPushed"
}

func (s ControlPushed) Description() string {
	return "ControlPushed - signal represents state of key with 'Control' type right off it was pressed on a device"
}
