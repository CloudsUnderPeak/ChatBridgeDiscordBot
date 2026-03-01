package gamble

import (
	"sync"
)

type Gamer struct {
	id    string
	name  string
	chips int64
	mu    sync.RWMutex
	dbMu  sync.RWMutex
}

const (
	GAME_BIGGER_NUMBER = "BiggerNumber"
	GAME_SLOT_MACHINE  = "SlotMachine"
)
