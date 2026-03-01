package ai

import (
	"discord-chatbot/discord/pkg/discordbot"
	"discord-chatbot/pkg/aiAgent"
	"discord-chatbot/pkg/logger"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var log = logger.GetLogger("api.ai")

func (a *apiInterface) Command(c *discordbot.Context) {
	var message string
	field := strings.Fields(c.Message.Content)
	if len(field) > 1 {
		username := c.Session.State.User.Username
		// c.Session.ChannelMessageSend(c.Message.ChannelID, "wait a minute...")
		c.Session.ChannelTyping(c.Message.ChannelID)
		id := fmt.Sprintf("%s-%s-%s", username, c.Message.Author.ID, c.Message.ChannelID)
		if _, ok := a.msgDb[id]; !ok {
			a.msgDb[id] = &aiAgent.MessageDataBase{}
		}

		command := strings.Join(field[1:], " ")
		bot, ok := a.bots[username]
		if !ok || bot == nil {
			message = "ai bot not initialized"
			c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
				Content:   message,
				Reference: &discordgo.MessageReference{MessageID: c.Message.ID, ChannelID: c.Message.ChannelID, GuildID: c.Message.GuildID},
			})
			return
		}
		msg, err := bot.CommandWithDatabase(a.msgDb[id], command)
		if err != nil {
			log.Debugf("| %15s | %-8s | %v",
				c.Message.ChannelID,
				c.Message.Author.ID,
				err,
			)
			message = fmt.Sprintf("ai command fail (%s)", err.Error())
		} else {
			var msgContents []string
			for _, m := range msg {
				msgContents = append(msgContents, m.Content)
			}
			message = strings.Join(msgContents, "\n")
		}
	} else {
		message = "can I help you?"
	}

	c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
		Content:   message,
		Reference: &discordgo.MessageReference{MessageID: c.Message.ID, ChannelID: c.Message.ChannelID, GuildID: c.Message.GuildID},
	})
}
