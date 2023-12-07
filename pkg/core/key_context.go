package core

import (
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"time"
)

type KeyContext struct {
	key      uint8
	velocity uint8
	usedAt   time.Time
	status   core.Signal
}
