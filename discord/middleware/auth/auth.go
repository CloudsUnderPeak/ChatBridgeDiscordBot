package auth

import (
	"discord-chatbot/discord/pkg/discordbot"
	pkgConfig "discord-chatbot/pkg/config"
	"discord-chatbot/pkg/logger"

	"golang.org/x/exp/slices"
)

var log = logger.GetLogger("auth")

type AuthFuncConfig struct {
	function string
	level    int
	bot      pkgConfig.BotConfig
	users    []pkgConfig.UserConfig
}

func AuthFunc(function string, level int) *AuthFuncConfig {
	return &AuthFuncConfig{
		function: function,
		level:    level,
	}
}

func (a *AuthFuncConfig) RegisterBot(bot pkgConfig.BotConfig) *AuthFuncConfig {
	a.bot = bot
	return a
}

func (a *AuthFuncConfig) RegisterUser(users []pkgConfig.UserConfig) *AuthFuncConfig {
	a.users = users
	return a
}

func (a *AuthFuncConfig) RegisterUserByConfig() *AuthFuncConfig {
	var users []pkgConfig.UserConfig
	for i := range pkgConfig.Users {
		user := *pkgConfig.Users[i]
		users = append(users, user)
	}
	a.users = users
	return a
}

func (a *AuthFuncConfig) Access(c *discordbot.Context) {
	pass := true
	for {
		if a.bot.Name != "" && a.bot.Token != "" {
			pass = a.verfiyBot(c)
			if !pass {
				log.Debugf("| %15s | %-8s | %-8s | Function Ignore",
					c.Message.ChannelID,
					c.Message.Author.ID,
					c.Message.Author.GlobalName,
				)
				break
			}
		}
		if len(a.users) > 0 {
			pass = a.verfiyUser(c)
			if !pass {
				log.Debugf("| %15s | %-8s | %-8s | User Level Ignore",
					c.Message.ChannelID,
					c.Message.Author.ID,
					c.Message.Author.GlobalName,
				)
				break
			}
		}
		break
	}
	if pass {
		c.Next()
	}
}

func (a *AuthFuncConfig) verfiyBot(c *discordbot.Context) bool {
	pass := true
	for {
		pass = slices.Contains(a.bot.Functions, a.function)
		if !pass && len(a.bot.Channels) > 0 {
			for _, channel := range a.bot.Channels {
				if c.Message.ChannelID == channel.Id {
					pass = slices.Contains(channel.Functions, a.function)
					break
				}
			}
		}
		break
	}
	return pass
}

func (a *AuthFuncConfig) verfiyUser(c *discordbot.Context) bool {
	pass := true
	for {
		if len(a.users) > 0 {
			mapUserConfigbyId := make(map[string]pkgConfig.UserConfig)
			mapUserConfigbyLevel := make(map[int][]pkgConfig.UserConfig)
			for _, user := range a.users {
				mapUserConfigbyId[user.Id] = user
				mapUserConfigbyLevel[user.Level] = append(mapUserConfigbyLevel[user.Level], user)
			}
			if len(mapUserConfigbyLevel[pkgConfig.UserLevel_Block]) > 0 {
				if user, ok := mapUserConfigbyId[c.Message.Author.ID]; ok {
					if user.Level == pkgConfig.UserLevel_Block {
						pass = false
						break
					}
				}
			}
			if len(mapUserConfigbyLevel[pkgConfig.UserLevel_Admin]) > 0 {
				userLevel := pkgConfig.UserLevel_Guest
				if user, ok := mapUserConfigbyId[c.Message.Author.ID]; ok {
					userLevel = user.Level
				}
				if userLevel < a.level {
					pass = false
					break
				}
			}
		}
		break
	}
	return pass
}
