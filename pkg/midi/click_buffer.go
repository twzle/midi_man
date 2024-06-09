package midi

type ClickBuffer map[uint8]*KeyContext

// Function returns context of key contained in click buffer by id
func (cb ClickBuffer) GetKeyContext(key uint8) (*KeyContext, bool) {
	val, ok := cb[key]
	return val, ok
}

// Function sets context of key contained in click buffer by id and current context
func (cb ClickBuffer) SetKeyContext(key uint8, midiKey KeyContext) {
	cb[key] = &midiKey
}
