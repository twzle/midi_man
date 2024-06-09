package model

// Representation of component click start event (NoteOn)
type NotePushed struct {
	Device    string `hubman:"device"`
	Namespace string `hubman:"namespace"`
	KeyCode   int    `hubman:"key_code"`
	Velocity  int    `hubman:"velocity"`
}

// Function returns string representation of model
func (s NotePushed) Code() string {
	return "NotePushed"
}

// Function returns string description of model
func (s NotePushed) Description() string {
	return "NotePushed - signal represents state of key with 'Note' type right off it was pressed on a device"
}

// Representation of component hold event
type NoteHold struct {
	Device    string `hubman:"device"`
	Namespace string `hubman:"namespace"`
	KeyCode   int    `hubman:"key_code"`
	Velocity  int    `hubman:"velocity"`
}

// Function returns string representation of model
func (s NoteHold) Code() string {
	return "NoteHold"
}

// Function returns string description of model
func (s NoteHold) Description() string {
	return "NoteHold - signal represents state of key with 'Note' type that is pressed for long"
}

// Representation of component click end event (NoteOff)
type NoteReleased struct {
	Device    string `hubman:"device"`
	Namespace string `hubman:"namespace"`
	KeyCode   int    `hubman:"key_code"`
	Velocity  int    `hubman:"velocity"`
}

// Function returns string representation of model
func (s NoteReleased) Code() string {
	return "NoteReleased"
}

// Function returns string description of model
func (s NoteReleased) Description() string {
	return "NoteReleased - signal represents state of key with 'Note' type right off it was released on a device"
}

// Representation of component click end event after hold
type NoteReleasedAfterHold struct {
	Device    string `hubman:"device"`
	Namespace string `hubman:"namespace"`
	KeyCode   int    `hubman:"key_code"`
	Velocity  int    `hubman:"velocity"`
}

// Function returns string representation of model
func (s NoteReleasedAfterHold) Code() string {
	return "NoteReleasedAfterHold"
}

// Function returns string description of model
func (s NoteReleasedAfterHold) Description() string {
	return "NoteReleasedAfterHold - signal represents state of key with 'Note' type right off it was released on a device after hold"
}

// ControlPushed В MIDI Control имеет только один тип событий "ControlChange",
// поэтому длительность и конец нажатия здесь не отслеживаются

// Representation of component click start event (ControlChanged)
type ControlPushed struct {
	Device    string `hubman:"device"`
	Namespace string `hubman:"namespace"`
	KeyCode   int    `hubman:"key_code"`
	Value     int    `hubman:"velocity"`
}

// Function returns string representation of model
func (s ControlPushed) Code() string {
	return "ControlPushed"
}

// Function returns string description of model
func (s ControlPushed) Description() string {
	return "ControlPushed - signal represents state of key with 'Control' type right off it was pressed on a device"
}


// Representation of namespace change event
type NamespaceChanged struct {
	Device       string `hubman:"device"`
	OldNamespace string `hubman:"old_namespace"`
	NewNamespace string `hubman:"new_namespace"`
}

// Function returns string representation of model
func (s NamespaceChanged) Code() string {
	return "NamespaceChanged"
}

// Function returns string description of model
func (s NamespaceChanged) Description() string {
	return "NamespaceChanged - signal represents successful namespace change"
}
