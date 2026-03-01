package ai

import (
	"discord-chatbot/pkg/aiAgent"
	pkgConfig "discord-chatbot/pkg/config"
	"sync"
)

var instance *apiInterface
var once sync.Once

type apiInterface struct {
	bots  map[string]*aiAgent.AiBot
	msgDb map[string]*aiAgent.MessageDataBase
}

func GetInstance() *apiInterface {
	once.Do(func() {
		instance = &apiInterface{}
		instance.Init()
	})
	return instance
}

func (a *apiInterface) Init() {
	a.bots = make(map[string]*aiAgent.AiBot)
	a.msgDb = make(map[string]*aiAgent.MessageDataBase)
}

func (a *apiInterface) InitBot(name string, agent pkgConfig.AiAgentConfig) error {
	apiKey := agent.ApiKey
	if apiKey == "" {
		switch agent.Provider {
		case aiAgent.ProviderOpenAI:
			apiKey = pkgConfig.ApiKey.OpenaiKey
		default:
			apiKey = pkgConfig.ApiKey.OpenaiKey
		}
		if apiKey == "" {
			return nil
		}
	}

	bot, err := aiAgent.NewAiBot(agent, apiKey)
	if err != nil {
		return err
	}
	bot.SetAiModel(agent.Provider, agent.Model, apiKey)
	bot.SetQueueLength(agent.QueueLength)
	a.bots[name] = bot
	a.bots[name].MessageDb.AddSystemMessages(agent.Prompt)
	return nil
}
