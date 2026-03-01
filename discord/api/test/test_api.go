package test

import (
	"discord-chatbot/discord/pkg/discordbot"
	"discord-chatbot/pkg/logger"
	"fmt"
)

var log = logger.GetLogger("test")

func (a *apiInterface) ErrorLog(c *discordbot.Context) {
	log.Errorf("Test error triggered by %s in channel %s", c.Message.Author.GlobalName, c.Message.ChannelID)
	c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("Error log sent."))
}
