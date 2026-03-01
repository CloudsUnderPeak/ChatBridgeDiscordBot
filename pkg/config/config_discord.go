package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type BotConfig struct {
	Name        string             `mapstructure:"name"`
	Token       string             `mapstructure:"token"`
	Enabled     bool               `mapstructure:"enabled"`
	Functions   []string           `mapstructure:"functions"`
	HelpUrl     string             `mapstructure:"helpUrl"`
	Channels    []ChannelConfig    `mapstructure:"channels"`
	LogChannels []LogChannelConfig `mapstructure:"logChannels"`
	AiAgent     AiAgentConfig      `mapstructure:"aiAgent"`
}

type ChannelConfig struct {
	Id        string   `mapstructure:"id"`
	Functions []string `mapstructure:"functions"`
}

type LogChannelConfig struct {
	Id string `mapstructure:"id"`
}

type UserConfig struct {
	Id    string `mapstructure:"id"`
	Level int    `mapstructure:"level"`
}

const (
	UserLevel_Block     = -1
	UserLevel_Guest     = 0
	UserLevel_User      = 1
	UserLevel_Moderator = 2
	UserLevel_Admin     = 3
)

type AiAgentConfig struct {
	Provider    string   `mapstructure:"provider"`
	ApiKey      string   `mapstructure:"apiKey"`
	Model       string   `mapstructure:"model"`
	QueueLength int      `mapstructure:"queueLength"`
	Prompt      []string `mapstructure:"prompt"`
}

type GameCenterConfig struct {
	GuessNumber GameGuessNumberConfig `mapstructure:"guessNumber"`
}

type GameGuessNumberConfig struct {
	Range int `mapstructure:"range"`
}

type GambleConfig struct {
	Principal    int64                    `mapstructure:"principal"`
	BiggerNumber GambleBiggerNumberConfig `mapstructure:"biggerNumber"`
	SlotMachine  GambleSlotMachineConfig  `mapstructure:"slotMachine"`
}

type GambleBiggerNumberConfig struct {
	Odds    float32 `mapstructure:"odds"`
	MinAnte int64   `mapstructure:"minAnte"`
}

type GambleSlotMachineConfig struct {
	MinAnte int64 `mapstructure:"minAnte"`
}

var (
	BotDefault = &BotConfig{
		Name:      "DEFAULT",
		Token:     "BOT_TOKEN",
		Enabled:   true,
		Functions: []string{"basic"},
		Channels:  []ChannelConfig{},
		AiAgent: AiAgentConfig{
			Provider:    "openai",
			ApiKey:      "",
			Model:       "gpt-4o-mini",
			QueueLength: 10,
			Prompt:      []string{},
		},
	}
	Bots = []*BotConfig{}

	UserDefault = &UserConfig{
		Id:    "CHANNEL_ID",
		Level: UserLevel_Admin,
	}
	Users = []*UserConfig{}

	GameCenter = &GameCenterConfig{
		GuessNumber: GameGuessNumberConfig{
			Range: 100,
		},
	}

	Gamble = &GambleConfig{
		Principal: 1000,
		BiggerNumber: GambleBiggerNumberConfig{
			Odds:    0.5,
			MinAnte: 100,
		},
		SlotMachine: GambleSlotMachineConfig{
			MinAnte: 100,
		},
	}
)

func initDiscord() {

	fmt.Println("Init Discord Config  ...")

	for {
		vp := viper.New()
		vp.SetConfigFile(Args.DiscordConfigPath)
		vp.SetConfigType("yaml")

		if err := vp.ReadInConfig(); err != nil {
			fmt.Printf("Config fail to read '%v': %v\n", Args.DiscordConfigPath, err)
			break
		}

		mapToList(vp, "bots", &Bots, BotDefault)
		mapToList(vp, "users", &Users, UserDefault)
		mapTo(vp, "gameCenter", &GameCenter)
		mapTo(vp, "gamble", &Gamble)

		initBots(vp, Bots)

		break
	}

	fmt.Println("Init Discord Config OK")
}

func initBots(vp *viper.Viper, bots []*BotConfig) {
	for i := range bots {
		botName := strings.ToUpper(bots[i].Name)
		botTokenVar := fmt.Sprintf("bots.%d.token", i)
		botTokenEnv := fmt.Sprintf("%s_BOT_TOKEN", botName)
		vp.BindEnv(botTokenVar, botTokenEnv)
	}
}
