package midi


type Control struct {
	Key              byte
	Rotate           bool
	ValueRange       [2]byte
	InitialValue     byte
	IncrementTrigger byte
	DecrementTrigger byte
}

func (md *MidiDevice) handleControls(controlKey byte, controlVelocity byte) (byte, bool) {
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
