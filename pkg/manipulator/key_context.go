package manipulator

import (
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"midi_manipulator/pkg/signals"
	"time"
)

type MidiKey struct {
	key      float64
	velocity float64
	usedAt   time.Time
	status   core.Signal
}

func (mk *MidiKey) getKeyCode() float64 {
	return mk.key
}

func (mk *MidiKey) getVelocity() float64 {
	return mk.velocity
}

func (mk *MidiKey) getUsedAt() time.Time {
	return mk.usedAt
}

func (mk *MidiKey) getStatus() core.Signal {
	return mk.status
}

func (mk *MidiKey) setStatus(status core.Signal) {
	mk.status = status
}

func (mk *MidiKey) isActive() bool {
	switch mk.getStatus().(type) {
	case nil, signals.NoteReleased, signals.ControlPushed:
		return false
	default:
		return true
	}
}

// KeyContext Можно сделать массив из контекстов, если нужно будет отслеживать несколько одновременных нажатий
type KeyContext struct {
	currentKey  MidiKey
	previousKey MidiKey
}

func (kctx *KeyContext) setCurrentKey(key MidiKey) {
	kctx.currentKey = key
}

func (kctx *KeyContext) getCurrentKey() MidiKey {
	return kctx.currentKey
}

func (kctx *KeyContext) setPreviousKey(key MidiKey) {
	kctx.previousKey = key
}

func (kctx *KeyContext) getPreviousKey() MidiKey {
	return kctx.previousKey
}

func (kctx *KeyContext) isPreviousKeyActive() bool {
	if &kctx.previousKey != nil {
		return kctx.previousKey.isActive()
	} else {
		return false
	}
}

func (kctx *KeyContext) compareKeys() bool {
	return kctx.previousKey == kctx.currentKey
}

func (kctx *KeyContext) compareStatuses() bool {
	return kctx.previousKey.getStatus() == kctx.currentKey.getStatus()
}
