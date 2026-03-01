package help

import (
	"discord-chatbot/discord/pkg/discordbot"
	pkgConfig "discord-chatbot/pkg/config"
	tr "discord-chatbot/pkg/translate"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/slices"
)

type helpContent struct {
	function    string
	name        string
	alias       []string
	command     string
	level       int
	description string
}

var userLevelText = map[int]string{
	pkgConfig.UserLevel_Guest:     "guest",
	pkgConfig.UserLevel_User:      "user",
	pkgConfig.UserLevel_Moderator: "moderator",
	pkgConfig.UserLevel_Admin:     "admin",
}

func GetContents() []helpContent {
	return []helpContent{
		{
			function:    "help",
			name:        "help",
			alias:       tr.Ts("discord.api.help.content.help.alias"),
			command:     tr.T("discord.api.help.content.help.command"),
			level:       0,
			description: tr.T("discord.api.help.content.help.desc"),
		},
		{
			function:    "basic",
			name:        "hi",
			alias:       tr.Ts("discord.api.help.content.hi.alias"),
			command:     tr.T("discord.api.help.content.hi.command"),
			level:       0,
			description: tr.T("discord.api.help.content.hi.desc"),
		},
		{
			function:    "ai",
			name:        "ai",
			alias:       tr.Ts("discord.api.help.content.ai.alias"),
			command:     tr.T("discord.api.help.content.ai.command"),
			level:       0,
			description: tr.T("discord.api.help.content.ai.desc"),
		},
		{
			function:    "gamecenter",
			name:        "guess",
			alias:       tr.Ts("discord.api.help.content.guess.alias"),
			command:     tr.T("discord.api.help.content.guess.command"),
			level:       0,
			description: tr.T("discord.api.help.content.guess.desc"),
		},
		{
			function:    "gamecenter",
			name:        "resetguess",
			alias:       tr.Ts("discord.api.help.content.resetguess.alias"),
			command:     tr.T("discord.api.help.content.resetguess.command"),
			level:       0,
			description: tr.T("discord.api.help.content.resetguess.desc"),
		},
		{
			function:    "gamecenter",
			name:        "1a2b",
			alias:       tr.Ts("discord.api.help.content.1a2b.alias"),
			command:     tr.T("discord.api.help.content.1a2b.command"),
			level:       0,
			description: tr.T("discord.api.help.content.1a2b.desc"),
		},
		{
			function:    "gamecenter",
			name:        "reset1a2b",
			alias:       tr.Ts("discord.api.help.content.reset1a2b.alias"),
			command:     tr.T("discord.api.help.content.reset1a2b.command"),
			level:       0,
			description: tr.T("discord.api.help.content.reset1a2b.desc"),
		},
		{
			function:    "gamble",
			name:        "rank",
			alias:       tr.Ts("discord.api.help.content.rank.alias"),
			command:     tr.T("discord.api.help.content.rank.command"),
			level:       0,
			description: tr.T("discord.api.help.content.rank.desc"),
		},
		{
			function:    "gamble",
			name:        "gamble",
			alias:       tr.Ts("discord.api.help.content.gamble.alias"),
			command:     tr.T("discord.api.help.content.gamble.command"),
			level:       0,
			description: tr.T("discord.api.help.content.gamble.desc"),
		},
		{
			function:    "gamble",
			name:        "slot",
			alias:       tr.Ts("discord.api.help.content.slot.alias"),
			command:     tr.T("discord.api.help.content.slot.command"),
			level:       0,
			description: tr.T("discord.api.help.content.slot.desc"),
		},
		{
			function:    "gamble",
			name:        "chips",
			alias:       tr.Ts("discord.api.help.content.chips.alias"),
			command:     tr.T("discord.api.help.content.chips.command"),
			level:       0,
			description: tr.T("discord.api.help.content.chips.desc"),
		},
		{
			function:    "gamble",
			name:        "repay",
			alias:       tr.Ts("discord.api.help.content.repay.alias"),
			command:     tr.T("discord.api.help.content.repay.command"),
			level:       0,
			description: tr.T("discord.api.help.content.repay.desc"),
		},
		{
			function:    "gamble",
			name:        "give",
			alias:       tr.Ts("discord.api.help.content.give.alias"),
			command:     tr.T("discord.api.help.content.give.command"),
			level:       3,
			description: tr.T("discord.api.help.content.give.desc"),
		},
	}
}

func GetAlias(name string) []string {
	for _, content := range GetContents() {
		if content.name == name {
			return content.alias
		}
	}
	return nil
}

func (a *apiInterface) Help(c *discordbot.Context) {
	if a.botConfig.HelpUrl != "" {
		c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
			Content:   tr.T("discord.api.help.url", a.botConfig.HelpUrl),
			Reference: &discordgo.MessageReference{MessageID: c.Message.ID, ChannelID: c.Message.ChannelID, GuildID: c.Message.GuildID},
		})
		return
	}

	var messages []string

	mapFunctions := make(map[string]bool)
	functions := a.botConfig.Functions
	for _, function := range functions {
		mapFunctions[function] = true
	}

	for _, channel := range a.botConfig.Channels {
		if channel.Id != c.Message.ChannelID {
			continue
		}
		for _, function := range channel.Functions {
			if _, ok := mapFunctions[function]; !ok {
				functions = append(functions, function)
				mapFunctions[function] = true
			}
		}
	}
	messages = append(messages, tr.T("discord.api.help.title"))
	for _, content := range GetContents() {
		if slices.Contains(functions, content.function) || content.function == "help" {
			messages = append(messages, fmt.Sprintf("** %s **   -  %s", content.command, userLevelText[content.level]))
			if len(content.alias) > 0 {
				messages = append(messages, fmt.Sprintf(" [ %s ]", strings.Join(content.alias, ", ")))
			}
			messages = append(messages, fmt.Sprintf("    %s\n", content.description))
		}
	}
	c.Session.ChannelMessageSend(c.Message.ChannelID, strings.Join(messages, "\n"))
}
