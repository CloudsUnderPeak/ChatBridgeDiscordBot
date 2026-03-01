package gamecenter

import (
	"sync"
)

type Gamer struct {
	id   string
	name string
}

type Game struct {
	id   string
	data interface{}
	lock sync.Mutex
}

const (
	GAME_GUESS_NUMBER   = "GuessNumber"
	GAME_BULLS_AND_COWS = "BullsAndCows" //1A2B
)
