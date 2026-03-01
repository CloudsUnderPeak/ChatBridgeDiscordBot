package basic

import (
	"discord-chatbot/discord/pkg/discordbot"
	tr "discord-chatbot/pkg/translate"

	"github.com/bwmarrin/discordgo"
)

func (a *apiInterface) Hi(c *discordbot.Context) {
	c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
		Content:   tr.T("discord.api.basic.hi"),
		Reference: &discordgo.MessageReference{MessageID: c.Message.ID, ChannelID: c.Message.ChannelID, GuildID: c.Message.GuildID},
	})
}
