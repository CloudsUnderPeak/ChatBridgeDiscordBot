package gamble

import (
	pkgConfig "discord-chatbot/pkg/config"
	tr "discord-chatbot/pkg/translate"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

func (g *Gamer) getBetNumerFromContent(content string, cell int) (int64, error) {

	if cell < 0 {
		return 0, fmt.Errorf("cell index out of range")
	}

	var bet int64
	parts := strings.Fields(content)

	invalidErr := fmt.Errorf(tr.T("discord.api.gamble.format_error"))
	if len(parts) > cell {
		if parts[cell] == "all" {
			bet = g.GetChips()
			return bet, nil
		}
		val, err := strconv.ParseInt(parts[cell], 10, 64)
		if err != nil || val < 0 {
			return 0, invalidErr
		}
		bet = val
	} else {
		return 0, invalidErr
	}

	return bet, nil
}

func biggerNumberGame(gamer *Gamer, bet int64) []string {

	var messages []string

	for {
		if gamer.GetChips() < bet {
			messages = append(messages, tr.T("discord.api.gamble.not_enough", gamer.name))
			break
		}

		minAnte := pkgConfig.Gamble.BiggerNumber.MinAnte
		if bet < minAnte {
			messages = append(messages, tr.T("discord.api.gamble.not_enough_min", minAnte))
			break
		}

		chance := rand.Float32()
		numbers := []int{
			rand.Intn(99),
			rand.Intn(99),
		}
		sort.Ints(numbers)
		if numbers[0] == numbers[len(numbers)-1] {
			numbers[len(numbers)-1] = numbers[len(numbers)-1] + 1
		}

		var playerNum, computerNum int
		if chance < pkgConfig.Gamble.BiggerNumber.Odds {
			playerNum = numbers[len(numbers)-1]
			computerNum = numbers[0]
			gamer.SetChipsByGame(bet, bet*2)
		} else {
			playerNum = numbers[0]
			computerNum = numbers[len(numbers)-1]
			gamer.SetChipsByGame(bet, 0)
		}

		messages = append(messages, tr.T("discord.api.gamble.bigger_number.show", gamer.name, playerNum, computerNum))
		if playerNum > computerNum {
			messages = append(messages, tr.T("discord.api.gamble.bigger_number.win", bet*2))
		} else {
			messages = append(messages, tr.T("discord.api.gamble.bigger_number.lose", bet))
		}

		break
	}

	return messages
}

func slotMachineGame(gamer *Gamer, bet int64) []string {
	var messages []string

	if gamer.GetChips() < bet {
		messages = append(messages, tr.T("discord.api.gamble.not_enough", gamer.name))
		return messages
	}

	minAnte := pkgConfig.Gamble.SlotMachine.MinAnte
	if bet < minAnte {
		messages = append(messages, tr.T("discord.api.gamble.not_enough_min", minAnte))
		return messages
	}

	symbols := tr.Ts("discord.api.gamble.slot_machine.symbols")
	reels := make([]string, 3)

	for i := range reels {
		idx := rand.Intn(len(symbols))
		reels[i] = symbols[idx]
	}

	count := make(map[string]int)
	for _, sym := range reels {
		count[sym]++
	}

	maxMatch := 0
	for _, cnt := range count {
		if cnt > maxMatch {
			maxMatch = cnt
		}
	}

	var multiplier float64
	switch maxMatch {
	case 3:
		multiplier = 7
	case 2:
		multiplier = 1.5
	default:
		multiplier = 0
	}

	gamer.SetChipsByGame(bet, int64(float64(bet)*multiplier))

	messages = append(messages, tr.T("discord.api.gamble.slot_machine.show", gamer.name, reels[0], reels[1], reels[2]))
	if maxMatch > 1 {
		messages = append(messages, tr.T("discord.api.gamble.slot_machine.win", bet, multiplier))
	} else {
		messages = append(messages, tr.T("discord.api.gamble.slot_machine.lose", bet))
	}

	return messages
}
