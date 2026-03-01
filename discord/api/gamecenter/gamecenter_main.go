package gamecenter

import (
	"sync"
)

var instance *apiInterface
var once sync.Once

type apiInterface struct {
	gamers map[string]*Gamer
	games  map[string]*Game
}

func GetInstance() *apiInterface {
	once.Do(func() {
		instance = &apiInterface{}
		instance.Init()
	})
	return instance
}

func (a *apiInterface) Init() {
	a.gamers = make(map[string]*Gamer)
	a.games = make(map[string]*Game)
}
