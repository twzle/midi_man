package manipulator

type ClickBuffer map[uint8]*KeyContext

func (cb *ClickBuffer) GetKeyContext(key uint8) (*KeyContext, bool) {
	val, ok := (*cb)[key]
	return val, ok
}

func (cb *ClickBuffer) SetKeyContext(key uint8, midiKey KeyContext) {
	(*cb)[key] = &midiKey
}
