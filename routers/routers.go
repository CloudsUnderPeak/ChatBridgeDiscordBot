package routers

import (
	"discord-chatbot/discord/api/ai"
	"discord-chatbot/discord/api/basic"
	"discord-chatbot/discord/api/gamble"
	"discord-chatbot/discord/api/gamecenter"
	"discord-chatbot/discord/api/help"
	"discord-chatbot/discord/api/test"
	authMw "discord-chatbot/discord/middleware/auth"
	"discord-chatbot/discord/pkg/discordbot"
	"discord-chatbot/discord/pkg/discordlogger"
	pkgConfig "discord-chatbot/pkg/config"
	"discord-chatbot/pkg/logger"
	"discord-chatbot/pkg/signal"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var RestartRouterErr error = errors.New("Restart")

var log = logger.GetLogger("routers")

var Bots map[string]*discordbot.Bot
var monitorBotAlive map[string]chan struct{}

func Run() error {

	for {

		Bots = make(map[string]*discordbot.Bot)
		monitorBotAlive = make(map[string]chan struct{})

		for i := range pkgConfig.Bots {
			botConfig := pkgConfig.Bots[i]

			token := botConfig.Token
			botToken := fmt.Sprintf("Bot %s", token)

			bot, _ := discordbot.NewBot(botConfig)

			Bots[botToken] = bot
		}

		for i := range Bots {
			bot := Bots[i]
			bot.Use(loggerMiddleware)
			bot.Handle("!help", help.GetAlias("help"), help.GetInstance(bot.GetConfig()).Help)
			mapFunctions := make(map[string]bool)
			for _, function := range bot.GetConfig().Functions {
				mapFunctions[function] = true
			}
			for _, channel := range bot.GetConfig().Channels {
				for _, function := range channel.Functions {
					mapFunctions[function] = true
				}
			}
			for function := range mapFunctions {
				switch function {
				case "basic":
					bot.Handle("!hi", help.GetAlias("hi"), authFuncMwGuest(function, *bot).Access, basic.GetInstance().Hi)

				case "gamecenter":
					bot.Handle("!guess", help.GetAlias("guess"), authFuncMwGuest(function, *bot).Access, gamecenter.GetInstance().Game(gamecenter.GAME_GUESS_NUMBER))
					bot.Handle("!resetguess", help.GetAlias("resetguess"), authFuncMwGuest(function, *bot).Access, gamecenter.GetInstance().ResetGame(gamecenter.GAME_GUESS_NUMBER))
					bot.Handle("!1a2b", help.GetAlias("1a2b"), authFuncMwGuest(function, *bot).Access, gamecenter.GetInstance().Game(gamecenter.GAME_BULLS_AND_COWS))
					bot.Handle("!reset1a2b", help.GetAlias("reset1a2b"), authFuncMwGuest(function, *bot).Access, gamecenter.GetInstance().ResetGame(gamecenter.GAME_BULLS_AND_COWS))
					bot.Handle("!peek1a2b", help.GetAlias("peek1a2b"), authFuncMwAdmin(function, *bot).Access, gamecenter.GetInstance().PeekGame(gamecenter.GAME_BULLS_AND_COWS))

				case "gamble":
					bot.Handle("!rank", help.GetAlias("rank"), authFuncMwGuest(function, *bot).Access, gamble.GetInstance().GetRankings)
					bot.Handle("!gamble", help.GetAlias("gamble"), authFuncMwGuest(function, *bot).Access, gamble.GetInstance().Game(gamble.GAME_BIGGER_NUMBER))
					bot.Handle("!slot", help.GetAlias("slot"), authFuncMwGuest(function, *bot).Access, gamble.GetInstance().Game(gamble.GAME_SLOT_MACHINE))
					bot.Handle("!chips", help.GetAlias("chips"), authFuncMwGuest(function, *bot).Access, gamble.GetInstance().GetChips)
					bot.Handle("!repay", help.GetAlias("repay"), authFuncMwGuest(function, *bot).Access, gamble.GetInstance().Repay)
					bot.Handle("!give", help.GetAlias("give"), authFuncMwGuest(function, *bot).Access, gamble.GetInstance().GiveChips)

				case "ai":
					bot.Handle("!ai", help.GetAlias("ai"), authFuncMwGuest(function, *bot).Access, ai.GetInstance().Command)

				case "test":
					bot.Handle("!testerror", nil, authFuncMwAdmin(function, *bot).Access, test.GetInstance().ErrorLog)
				}
			}
		}

		for i := range Bots {
			bot := Bots[i]
			err := bot.Start()
			if err != nil {
				log.Warnf("The %s connect to Discord fail: %v", bot.GetName(), err)
				bot.SetEnabled(false)
				continue
			}
			// Register Discord log hook if logChannels configured
			if len(bot.GetConfig().LogChannels) > 0 {
				hook := discordlogger.NewDiscordHook(bot.GetSession(), bot.GetConfig().LogChannels)
				logrus.AddHook(hook)
			}
			// InitBot after Start() so Session.State.User.Username is available
			if err := ai.GetInstance().InitBot(bot.GetUsername(), bot.GetConfig().AiAgent); err != nil {
				log.Warnf("The %s init AI agent fail: %v", bot.GetName(), err)
			}
			stopChan := make(chan struct{})
			monitorBotAlive[bot.GetTokenWithPrefix()] = stopChan
			go monitorBot(bot, stopChan)
			defer bot.Stop()
		}

		quit := make(chan os.Signal)
		signal.RouterSignalNotify(quit)
		signalValue := <-quit

		for _, ch := range monitorBotAlive {
			close(ch)
		}

		log.Info("Bot Exiting")

		restartMainSignal, getRestartMainSignalErr := signal.GetRestartSignal("main")
		if signalValue == restartMainSignal && getRestartMainSignalErr == nil {
			return RestartRouterErr
		}

		log.Info("Bot Shutdown")
		break
	}

	return nil
}

func loggerMiddleware(c *discordbot.Context) {
	log.Infof(" | %15s | %-8s | %-8s ",
		c.Message.ChannelID,
		c.Message.Author.GlobalName,
		c.Message.Content,
	)
	c.Next()
}

func authFuncMwGuest(function string, bot discordbot.Bot) *authMw.AuthFuncConfig {
	return authMw.AuthFunc(function, pkgConfig.UserLevel_Guest).RegisterBot(bot.GetConfig()).RegisterUserByConfig()
}

func authFuncMwAdmin(function string, bot discordbot.Bot) *authMw.AuthFuncConfig {
	return authMw.AuthFunc(function, pkgConfig.UserLevel_Admin).RegisterBot(bot.GetConfig()).RegisterUserByConfig()
}

func monitorBot(bot *discordbot.Bot, stop <-chan struct{}) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	retryCount := 0
	maxRetryCount := 30
	for {
		select {
		case <-ticker.C:
			if !bot.IsConnected() {
				retryCount++
				// log.Debugf("Bot %s connect fail, count %d", bot.GetName(), retryCount)
				if retryCount > maxRetryCount {
					log.Warnf("Bot %s connect fail, reconnecting", bot.GetName())
					bot.Stop()
					if err := bot.Start(); err != nil {
						log.Errorf("Bot %s connect fail: %v", bot.GetName(), err)
						retryCount = 0
					} else {
						log.Infof("Bot %s connect success", bot.GetName())
					}
				}
			} else {
				retryCount = 0
			}
		case <-stop:
			log.Infof("Stop monitor Bot %s", bot.GetName())
			return
		}
	}
}
