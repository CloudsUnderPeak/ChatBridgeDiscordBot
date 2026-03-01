package discordbot

import (
	pkgConfig "discord-chatbot/pkg/config"
	"discord-chatbot/pkg/logger"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var log = logger.GetLogger("discordbot")

func NewBot(config *pkgConfig.BotConfig) (*Bot, error) {

	token := fmt.Sprintf("Bot %s", config.Token)

	discord, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}
	bot := &Bot{
		config:          config,
		session:         discord,
		middlewares:     []HandlerFunc{},
		prefixHandlers:  make(map[string][]HandlerFunc),
		keywordContexts: []KeywordContext{},
	}
	discord.AddHandler(bot.contextHandler)
	return bot, nil
}

func (b *Bot) Use(mw HandlerFunc) {
	b.middlewares = append(b.middlewares, mw)
}

func (b *Bot) Handle(cmd string, alias []string, handlers ...HandlerFunc) {
	var cmds []string
	cmds = append(cmds, cmd)
	cmds = append(cmds, alias...)
	for _, c := range cmds {
		b.prefixHandlers[c] = handlers
	}
}

func (b *Bot) HandleKeyword(keywords []string, h HandlerFunc) {
	b.keywordContexts = append(b.keywordContexts, KeywordContext{keywords, h})
}

func (c *Context) Next() {
	c.index++
	if c.index < len(c.handlers) {
		c.handlers[c.index](c)
	}
}

func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return err
	}
	log.Infof("Bot %s is now running", b.config.Name)
	return nil
}

func (b *Bot) Stop() error {
	log.Infof("Bot %s is now stopping", b.config.Name)
	return b.session.Close()
}

func (b *Bot) contextHandler(session *discordgo.Session, message *discordgo.MessageCreate) {

	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == session.State.User.ID {
		return
	}

	pipeline := make([]HandlerFunc, 0, len(b.middlewares)+1)
	pipeline = append(pipeline, b.middlewares...)

	existHandler := false
	field := strings.Fields(message.Content)
	if len(field) > 0 {
		prefix := field[0]
		prefix = strings.ReplaceAll(prefix, "！", "!")
		if handlers, ok := b.prefixHandlers[prefix]; ok {
			pipeline = append(pipeline, handlers...)
			existHandler = true
		}
	}
	if !existHandler {
		for _, keywordContext := range b.keywordContexts {
			for _, keyword := range keywordContext.keywords {
				if strings.Contains(message.Content, keyword) {
					pipeline = append(pipeline, keywordContext.handler)
					break
				}
			}
		}
	}

	ctx := &Context{
		Session:  session,
		Message:  message,
		handlers: pipeline,
		index:    0,
	}
	ctx.handlers[0](ctx)
}

func (b *Bot) GetConfig() pkgConfig.BotConfig {
	return *b.config
}

func (b *Bot) GetName() string {
	return b.config.Name
}

func (b *Bot) GetUsername() string {
	return b.session.State.User.Username
}

func (b *Bot) GetEnabled() bool {
	return b.config.Enabled
}

func (b *Bot) SetEnabled(enabled bool) bool {
	b.config.Enabled = enabled
	return b.config.Enabled
}

func (b *Bot) GetToken() string {
	return b.config.Token
}

func (b *Bot) GetTokenWithPrefix() string {
	return fmt.Sprintf("Bot %s", b.config.Token)
}

func (b *Bot) GetSession() *discordgo.Session {
	return b.session
}

func (b *Bot) IsConnected() bool {
	if b.session.HeartbeatLatency() > 0 {
		return true
	}
	return false
}
