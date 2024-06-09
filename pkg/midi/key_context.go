package midi

import (
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"time"
)

// Representation of key context entity
type KeyContext struct {
	key      uint8
	velocity uint8
	usedAt   time.Time
	status   core.Signal
}
