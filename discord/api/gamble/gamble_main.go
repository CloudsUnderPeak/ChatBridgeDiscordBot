package gamble

import (
	"sync"
)

var instance *apiInterface
var once sync.Once

type apiInterface struct {
	gamers map[string]*Gamer
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
	a.InitGamble()
}
