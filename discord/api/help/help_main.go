package help

import (
	pkgConfig "discord-chatbot/pkg/config"
	"sync"
)

var instance map[string]*apiInterface
var instanceOnce map[string]*sync.Once
var once sync.Once

type apiInterface struct {
	botConfig pkgConfig.BotConfig
}

func GetInstance(bot pkgConfig.BotConfig) *apiInterface {
	once.Do(func() {
		instance = make(map[string]*apiInterface)
		instanceOnce = make(map[string]*sync.Once)
	})
	if _, ok := instanceOnce[bot.Token]; !ok {
		instanceOnce[bot.Token] = &sync.Once{}
		instanceOnce[bot.Token].Do(func() {
			instance[bot.Token] = &apiInterface{botConfig: bot}
			instance[bot.Token].Init()
		})
	}
	return instance[bot.Token]
}

func (a *apiInterface) Init() {
}
