package manipulator

import (
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"midi_manipulator/pkg/signals"
	"time"
)

type KeyContext struct {
	key      int
	velocity int
	usedAt   time.Time
	status   core.Signal
}

func (kctx *KeyContext) getKeyCode() int {
	return kctx.key
}

func (kctx *KeyContext) getVelocity() int {
	return kctx.velocity
}

func (kctx *KeyContext) getUsedAt() time.Time {
	return kctx.usedAt
}

func (kctx *KeyContext) getStatus() core.Signal {
	return kctx.status
}

func (kctx *KeyContext) setStatus(status core.Signal) {
	kctx.status = status
}

func (kctx *KeyContext) isActive() bool {
	switch kctx.getStatus().(type) {
	case nil, signals.NoteReleased, signals.ControlPushed:
		return false
	default:
		return true
	}
}
