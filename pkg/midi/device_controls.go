package midi

// Representation of control entity
type Control struct {
	Key              int
	Rotate           bool
	ValueRange       [2]int
	InitialValue     int
	IncrementTrigger int
	DecrementTrigger int
}

// Function handles behaviour of control by id and velocity and returns modified values
func (md *MidiDevice) handleControls(controlKey int, controlVelocity int) (int, bool) {
	control, ok := md.controls[controlKey]
	if ok && !control.Rotate {
		if control.IncrementTrigger == controlVelocity && control.InitialValue < control.ValueRange[1]{
			control.InitialValue++
			return control.InitialValue, true
		} else if control.DecrementTrigger == controlVelocity && control.InitialValue > control.ValueRange[0] {
			control.InitialValue--
			return control.InitialValue, true
		} else {
			return control.InitialValue, false // unmodified value banned 
		}
	} else {
		return controlVelocity, true // unfiltered value accepted
	}
}
