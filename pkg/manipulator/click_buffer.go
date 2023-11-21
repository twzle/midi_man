package manipulator

type ClickBuffer struct {
	buffer map[int]KeyContext
}

func (cb *ClickBuffer) getKeyContext(key int) (KeyContext, bool) {
	val, ok := cb.buffer[key]
	return val, ok
}

func (cb *ClickBuffer) setKeyContext(key int, midiKey KeyContext) {
	cb.buffer[key] = midiKey
}
