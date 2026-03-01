package discordbot

import (
	pkgConfig "discord-chatbot/pkg/config"

	"github.com/bwmarrin/discordgo"
)

type HandlerFunc func(*Context)

type KeywordContext struct {
	keywords []string
	handler  HandlerFunc
}

type Bot struct {
	session         *discordgo.Session
	config          *pkgConfig.BotConfig
	middlewares     []HandlerFunc
	prefixHandlers  map[string][]HandlerFunc
	keywordContexts []KeywordContext
}

type Context struct {
	Session  *discordgo.Session
	Message  *discordgo.MessageCreate
	handlers []HandlerFunc
	index    int
}
