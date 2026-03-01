package gamecenter

import (
	"discord-chatbot/discord/pkg/discordbot"
	tr "discord-chatbot/pkg/translate"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (a *apiInterface) Game(game string) func(c *discordbot.Context) {
	return func(c *discordbot.Context) {
		var messages []string
		content := c.Message.Content
		id := c.Message.Author.ID
		name := c.Message.Author.GlobalName
		gamer := a.RegistGamer(id, name)

		switch game {
		case GAME_GUESS_NUMBER:
			var number int
			parts := strings.Fields(content)
			if len(parts) > 1 {
				val, err := strconv.Atoi(parts[1])
				if err != nil {
					number = 0
				} else {
					number = val
				}
			} else {
				number = 0
			}
			gameId := fmt.Sprintf("%s-%s", game, c.Message.ChannelID)
			if _, ok := a.games[gameId]; !ok {
				a.games[gameId] = &Game{
					id: gameId,
				}
			}
			msgs := a.games[gameId].guessNumberGame(*gamer, number)
			messages = append(messages, msgs...)

		case GAME_BULLS_AND_COWS:
			var numberStr string
			parts := strings.Fields(content)
			if len(parts) > 1 {
				numberStr = parts[1]
			}
			gameId := fmt.Sprintf("%s-%s", game, c.Message.ChannelID)
			if _, ok := a.games[gameId]; !ok {
				a.games[gameId] = &Game{
					id: gameId,
				}
			}
			msgs := a.games[gameId].bullsAndCowsGame(*gamer, numberStr)
			messages = append(messages, msgs...)

		default:
			msg := tr.T("discord.api.gamecenter.not_available", game)
			messages = append(messages, msg)
		}
		c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
			Content:   strings.Join(messages, "\n"),
			Reference: &discordgo.MessageReference{MessageID: c.Message.ID, ChannelID: c.Message.ChannelID, GuildID: c.Message.GuildID},
		})
	}
}

func (a *apiInterface) ResetGame(game string) func(c *discordbot.Context) {
	return func(c *discordbot.Context) {
		var messages []string

		switch game {
		case GAME_GUESS_NUMBER:
			fallthrough
		case GAME_BULLS_AND_COWS:
			gameId := fmt.Sprintf("%s-%s", game, c.Message.ChannelID)
			if _, ok := a.games[gameId]; !ok {
				a.games[gameId] = &Game{
					id: gameId,
				}
			}
			a.games[gameId].lock.Lock()
			defer a.games[gameId].lock.Unlock()
			a.games[gameId].data = nil

			messages = append(messages, tr.T("discord.api.gamecenter.reset_success"))
		}
		c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
			Content:   strings.Join(messages, "\n"),
			Reference: &discordgo.MessageReference{MessageID: c.Message.ID, ChannelID: c.Message.ChannelID, GuildID: c.Message.GuildID},
		})
	}
}

func (a *apiInterface) PeekGame(game string) func(c *discordbot.Context) {
	return func(c *discordbot.Context) {
		var messages []string
		switch game {
		case GAME_BULLS_AND_COWS:
			gameId := fmt.Sprintf("%s-%s", game, c.Message.ChannelID)
			g, ok := a.games[gameId]
			if !ok {
				messages = append(messages, tr.T("discord.api.gamecenter.bulls_and_cows.peek_no_game"))
			} else {
				g.lock.Lock()
				defer g.lock.Unlock()
				digits, ok := g.data.([]int)
				if !ok || len(digits) < 4 {
					messages = append(messages, tr.T("discord.api.gamecenter.bulls_and_cows.peek_no_game"))
				} else {
					answer := fmt.Sprintf("%d%d%d%d", digits[0], digits[1], digits[2], digits[3])
					messages = append(messages, tr.T("discord.api.gamecenter.bulls_and_cows.peek", answer))
				}
			}
		}
		c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
			Content:   strings.Join(messages, "\n"),
			Reference: &discordgo.MessageReference{MessageID: c.Message.ID, ChannelID: c.Message.ChannelID, GuildID: c.Message.GuildID},
		})
	}
}

func (a *apiInterface) RegistGamer(id string, name string) *Gamer {
	if _, ok := a.gamers[id]; !ok {
		gamer := &Gamer{
			id:   id,
			name: name,
		}
		a.gamers[id] = gamer
	} else {
		a.gamers[id].name = name
	}

	return a.gamers[id]
}
