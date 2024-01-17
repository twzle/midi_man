package midi


type Control struct {
	Key              uint8
	Rotate           bool
	ValueRange       [2]uint8
	InitialValue     uint8
	IncrementTrigger uint8
	DecrementTrigger uint8
}

func (md *MidiDevice) handleControls(controlKey uint8, controlVelocity uint8) (uint8, bool) {
	for _, control := range md.controls {
		if control.Key == controlKey && !control.Rotate {
			if control.IncrementTrigger == controlVelocity && control.InitialValue < control.ValueRange[1]{
				control.InitialValue++
				return control.InitialValue, true
			} else if control.DecrementTrigger == controlVelocity && control.InitialValue > control.ValueRange[0] {
				control.InitialValue--
				return control.InitialValue, true
			} else {
				return control.InitialValue, false // unmodified value banned
			}
		}
	}
	return controlVelocity, true // unfiltered value accepted
}
