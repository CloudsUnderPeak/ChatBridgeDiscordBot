package gamble

import (
	"discord-chatbot/discord/pkg/discordbot"
	pkgConfig "discord-chatbot/pkg/config"
	"discord-chatbot/pkg/logger"
	tr "discord-chatbot/pkg/translate"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var log = logger.GetLogger("api.gamble")

func (a *apiInterface) InitGamble() {
	gamers, err := ReadAllDbGamer()
	if err != nil {
		log.Warnf("InitGamble error: %v", err)
	} else {
		for _, gamer := range gamers {
			a.gamers[gamer.id] = gamer
		}
	}
}

func (a *apiInterface) Game(game string) func(c *discordbot.Context) {
	return func(c *discordbot.Context) {
		var messages []string
		content := c.Message.Content
		id := c.Message.Author.ID
		name := c.Message.Author.GlobalName
		gamer := a.RegistGamer(id, name)

		switch game {
		case GAME_BIGGER_NUMBER:
			bet, err := gamer.getBetNumerFromContent(content, 1)
			if err != nil {
				messages = append(messages, err.Error())
			} else {
				msgs := biggerNumberGame(gamer, bet)
				messages = append(messages, msgs...)
			}
		case GAME_SLOT_MACHINE:
			bet, err := gamer.getBetNumerFromContent(content, 1)
			if err != nil {
				messages = append(messages, err.Error())
			} else {
				msgs := slotMachineGame(gamer, bet)
				messages = append(messages, msgs...)
			}
		default:
			msg := tr.T("discord.api.gamble.not_available", game)
			messages = append(messages, msg)
		}
		c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
			Content:   strings.Join(messages, "\n"),
			Reference: &discordgo.MessageReference{MessageID: c.Message.ID, ChannelID: c.Message.ChannelID, GuildID: c.Message.GuildID},
		})
	}
}

func (a *apiInterface) GetChips(c *discordbot.Context) {
	content := c.Message.Content
	parts := strings.Fields(content)
	var message string
	if len(parts) > 1 {
		for _, gamer := range a.gamers {
			if gamer.name == parts[1] {
				message = tr.T("discord.api.gamble.get_chips", gamer.name, gamer.GetChips())
			}
		}
	} else {
		id := c.Message.Author.ID
		name := c.Message.Author.GlobalName
		gamer := a.RegistGamer(id, name)
		chips := gamer.GetChips()
		message = tr.T("discord.api.gamble.get_chips", name, chips)
	}
	c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
		Content:   message,
		Reference: &discordgo.MessageReference{MessageID: c.Message.ID, ChannelID: c.Message.ChannelID, GuildID: c.Message.GuildID},
	})
}

func (a *apiInterface) GiveChips(c *discordbot.Context) {
	content := c.Message.Content
	parts := strings.Fields(content)
	if len(parts) > 2 {
		isAll := false
		bet, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil || bet <= 0 {
			return
		}
		if parts[1] == "all" || parts[1] == "All" {
			isAll = true
		}
		for _, gamer := range a.gamers {
			if isAll || gamer.name == parts[1] {
				gamer.SetChips(gamer.GetChips() + bet)
				var messages []string
				messages = append(messages, tr.T("discord.api.gamble.give_chips", bet))
				messages = append(messages, tr.T("discord.api.gamble.get_chips", gamer.name, gamer.GetChips()))
				c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
					Content:   strings.Join(messages, "\n"),
					Reference: &discordgo.MessageReference{MessageID: c.Message.ID, ChannelID: c.Message.ChannelID, GuildID: c.Message.GuildID},
				})
			}
		}
	}
}

func (a *apiInterface) Repay(c *discordbot.Context) {
	id := c.Message.Author.ID
	name := c.Message.Author.GlobalName
	gamer := a.RegistGamer(id, name)
	chips := gamer.GetChips()
	if chips == 0 {
		gamer.SetChips(pkgConfig.Gamble.Principal)
	}
	message := tr.T("discord.api.gamble.get_chips", name, chips)
	c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
		Content:   message,
		Reference: &discordgo.MessageReference{MessageID: c.Message.ID, ChannelID: c.Message.ChannelID, GuildID: c.Message.GuildID},
	})
}

func (a *apiInterface) GetRankings(c *discordbot.Context) {
	gamersList := make([]*Gamer, 0, len(a.gamers))
	for _, g := range a.gamers {
		gamersList = append(gamersList, g)
	}
	sort.Slice(gamersList, func(i, j int) bool {
		return gamersList[i].chips > gamersList[j].chips
	})
	if len(gamersList) > 3 {
		gamersList = gamersList[:3]
	}
	var messages []string
	messages = append(messages, tr.T("discord.api.gamble.top3"))
	for i, top := range gamersList {
		messages = append(messages, tr.T("discord.api.gamble.rank_number", i+1, top.name, top.chips))
	}
	c.Session.ChannelMessageSendComplex(c.Message.ChannelID, &discordgo.MessageSend{
		Content:   strings.Join(messages, "\n"),
		Reference: &discordgo.MessageReference{MessageID: c.Message.ID, ChannelID: c.Message.ChannelID, GuildID: c.Message.GuildID},
	})
}

func (a *apiInterface) RegistGamer(id string, name string) *Gamer {
	if _, ok := a.gamers[id]; !ok {
		gamer := &Gamer{
			id:    id,
			name:  name,
			chips: pkgConfig.Gamble.Principal,
			mu:    sync.RWMutex{},
			dbMu:  sync.RWMutex{},
		}
		a.gamers[id] = gamer
		CreateDbGamer(gamer)
	} else {
		if a.gamers[id].name != name {
			a.gamers[id].name = name
			a.gamers[id].UpdateDbName()
		}
	}

	return a.gamers[id]
}

func (g *Gamer) GetChips() int64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.chips
}

func (g *Gamer) SetChips(chips int64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.chips = chips
	g.UpdateDbChips()
}

func (g *Gamer) SetChipsByGame(ante int64, bonus int64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.chips = g.chips - ante + bonus
	g.UpdateDbChips()
}
