package signals

type NotePushed struct {
	KeyCode  float64 `hubman:"key_code"`
	Velocity float64 `hubman:"velocity"`
}

func (c NotePushed) Code() string {
	return "NotePushed"
}

func (c NotePushed) Description() string {
	return "NotePushed - signal represents state of key with 'Note' type right off it was pressed on a device"
}

type NoteHold struct {
	KeyCode  float64 `hubman:"key_code"`
	Velocity float64 `hubman:"velocity"`
}

func (c NoteHold) Code() string {
	return "NoteHold"
}

func (c NoteHold) Description() string {
	return "NoteHold - signal represents state of key with 'Note' type that is pressed for long"
}

type NoteReleased struct {
	KeyCode  float64 `hubman:"key_code"`
	Velocity float64 `hubman:"velocity"`
}

func (c NoteReleased) Code() string {
	return "NoteReleased"
}

func (c NoteReleased) Description() string {
	return "NoteReleased - signal represents state of key with 'Note' type right off it was released on a device"
}

// ControlPushed В MIDI Control имеет только один тип событий "ControlChange",
// поэтому длительность и конец нажатия здесь не отслеживаются
type ControlPushed struct {
	KeyCode float64 `hubman:"key_code"`
	Value   float64 `hubman:"velocity"`
}

func (c ControlPushed) Code() string {
	return "ControlPushed"
}

func (c ControlPushed) Description() string {
	return "ControlPushed - signal represents state of key with 'Control' type right off it was pressed on a device"
}
