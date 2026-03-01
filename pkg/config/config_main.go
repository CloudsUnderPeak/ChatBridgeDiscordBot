package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type ArgsConfig struct {
	ConfigPath          string
	DiscordConfigPath   string
	DiscordDatabasePath string
	TranslationPath     string
	ResourcePath        string
}

type LogConfig struct {
	Path  string `mapstructure:"path"`
	Level string `mapstructure:"level"`
}

type ApiKeyConfig struct {
	OpenaiKey string `mapstructure:"openaiKey"`
}

type aiAgentConfig struct {
	Token    string `mapstructure:"token"`
	Provider string `mapstructure:"provider"`
	Model    string `mapstructure:"model"`
}

type SystemConfig struct {
	Language string `mapstructure:"language"`
}

var Args = &ArgsConfig{}

var (
	Log = &LogConfig{
		Path:  "/var/log/chatbot.log",
		Level: "debug",
	}

	ApiKey = &ApiKeyConfig{
		OpenaiKey: "",
	}

	System = &SystemConfig{
		Language: "zh",
	}
)

var rootDir string

func init() {

	fmt.Println("Init Config  ...")

	for {
		pwdPath, err := os.Getwd()
		if err != nil {
			fmt.Printf("Get pwd path fail: %v\n", err)
			pwdPath = os.Args[0]
		}

		rootDir, filepathError := filepath.Abs(pwdPath)
		if filepathError != nil {
			fmt.Printf("Config path error:%v\n", filepathError)
		}

		Args.ConfigPath = filepath.Join(rootDir, "conf/app.yaml")
		Args.DiscordConfigPath = filepath.Join(rootDir, "conf/discord.yaml")
		Args.DiscordDatabasePath = filepath.Join(rootDir, "data/discord.db")
		Args.ResourcePath = filepath.Join(rootDir, "resource")
		Args.TranslationPath = filepath.Join(rootDir, "conf/translations.json")

		vp := viper.New()
		vp.SetConfigFile(Args.ConfigPath)
		vp.SetConfigType("yaml")

		if err := vp.ReadInConfig(); err != nil {
			fmt.Printf("Config fail to read '%v': %v\n", Args.ConfigPath, err)
			break
		}

		mapTo(vp, "log", &Log)
		mapTo(vp, "apiKey", &ApiKey)
		mapTo(vp, "system", &System)

		initDiscord()

		if ApiKey.OpenaiKey == "" {
			vp.BindEnv("apiKey.openaiKey", "OPENAI_API_KEY")
			ApiKey.OpenaiKey = vp.GetString("apiKey.openaiKey")
		}

		break
	}

	fmt.Println("Init Config OK")
}
