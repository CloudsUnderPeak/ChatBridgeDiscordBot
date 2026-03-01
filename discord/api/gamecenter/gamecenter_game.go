package gamecenter

import (
	pkgConfig "discord-chatbot/pkg/config"
	tr "discord-chatbot/pkg/translate"
	"math/rand"
)

func (g *Game) guessNumberGame(gamer Gamer, number int) []string {

	var messages []string

	for {
		g.lock.Lock()
		defer g.lock.Unlock()

		if g.data == nil {
			g.data = 0
		} else if _, ok := g.data.(int); !ok {
			g.data = 0
		}

		if g.data.(int) == 0 {
			messages = append(messages, tr.T("discord.api.gamecenter.guess_number.welcome"))
			messages = append(messages, tr.T("discord.api.gamecenter.guess_number.input", pkgConfig.GameCenter.GuessNumber.Range))
			messages = append(messages, tr.T("discord.api.gamecenter.guess_number.gameplay"))
			g.data = rand.Intn(pkgConfig.GameCenter.GuessNumber.Range) + 1
		}

		if number < 1 || number > pkgConfig.GameCenter.GuessNumber.Range {
			messages = append(messages, tr.T("discord.api.gamecenter.guess_number.format_error", pkgConfig.GameCenter.GuessNumber.Range))
			break
		}

		secretNumber := g.data.(int)
		if number == secretNumber {
			messages = append(messages, tr.T("discord.api.gamecenter.guess_number.congratulations", gamer.name, g.data.(int)))
			g.data = 0
		} else if number < secretNumber {
			messages = append(messages, tr.T("discord.api.gamecenter.guess_number.lower_number", gamer.name, number))
		} else {
			messages = append(messages, tr.T("discord.api.gamecenter.guess_number.higher_number", gamer.name, number))
		}

		break
	}

	return messages
}

func (g *Game) bullsAndCowsGame(gamer Gamer, numberStr string) []string {
	var messages []string

	for {
		g.lock.Lock()
		defer g.lock.Unlock()

		if g.data == nil {
			g.data = []int{}
		} else if _, ok := g.data.([]int); !ok {
			g.data = []int{}
		}

		if len(g.data.([]int)) < 4 {
			messages = append(messages, tr.T("discord.api.gamecenter.bulls_and_cows.welcome"))
			messages = append(messages, tr.T("discord.api.gamecenter.bulls_and_cows.input"))
			messages = append(messages, tr.T("discord.api.gamecenter.bulls_and_cows.gameplay"))
			digits := rand.Perm(10)
			g.data = digits[:4]
		}

		guess := make([]int, 4)
		used := make(map[rune]bool)
		valid := true
		if len(numberStr) != 4 {
			valid = false
		} else {
			for i, ch := range numberStr {
				if ch < '0' || ch > '9' || used[ch] {
					valid = false
					break
				}
				used[ch] = true
				guess[i] = int(ch - '0')
			}
		}
		if !valid {
			messages = append(messages, tr.T("discord.api.gamecenter.bulls_and_cows.format_error"))
			break
		}

		a, b := evaluateBullsAndCows(g.data.([]int), guess)

		messages = append(messages, tr.T("discord.api.gamecenter.bulls_and_cows.user_guess", gamer.name, numberStr, a, b))
		if a == 4 {
			messages = append(messages, tr.T("discord.api.gamecenter.bulls_and_cows.congratulations", gamer.name))
			g.data = []int{}
		}

		break
	}

	return messages
}

func evaluateBullsAndCows(secret []int, guess []int) (int, int) {
	a, b := 0, 0
	for i := 0; i < 4; i++ {
		if guess[i] == secret[i] {
			a++
		} else {
			for j := 0; j < 4; j++ {
				if guess[i] == secret[j] {
					b++
				}
			}
		}
	}
	return a, b
}
