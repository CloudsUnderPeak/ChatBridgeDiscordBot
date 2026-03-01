package main

import (
	pkgConfig "discord-chatbot/pkg/config"
	"discord-chatbot/pkg/logger"
	pkgSql "discord-chatbot/pkg/sql"
	"discord-chatbot/pkg/translate"
	"discord-chatbot/routers"
	"io"
	"os"
	"runtime"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"

	"github.com/sirupsen/logrus"
)

var (
	VERSION string
	DATE    string
	COMMIT  string
	ARCH    string
)

var log = logger.GetLogger("main")

func main() {

	for {
		if len(VERSION) == 0 {
			VERSION = "?.?.?"
		}
		if len(DATE) == 0 {
			DATE = "?"
		}
		ARCH = runtime.GOARCH

		initLog()

		log.Infof("Version:%s Build:%s Commit:%s Arch:%s", VERSION, DATE, COMMIT, ARCH)

		if err := translate.InitTranslate(); err != nil {
			log.Errorf("Init Translate fail: %v", err)
			break
		}
		translate.SetLang(pkgConfig.System.Language)
		log.Infof("Language: %s", pkgConfig.System.Language)

		if err := pkgSql.InitSql(); err != nil {
			log.Errorf("Init Sql fail: %v", err)
			break
		}

		// bot run
		routerRunErr := routers.Run()

		if routerRunErr == routers.RestartRouterErr {
			log.Info("Bot Restart")
			continue
		}

		break
	}
}

func initLog() {
	if pkgConfig.Log.Level == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Set custom formatter
	logrus.SetFormatter(&logger.PackageFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Register caller hook (automatically adds function name and file:line)
	logrus.AddHook(logger.NewCallerHook(false))

	if len(pkgConfig.Log.Path) > 0 {
		file, err := os.OpenFile(pkgConfig.Log.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err == nil {
			multiWriter := io.MultiWriter(os.Stdout, file)
			logrus.SetOutput(multiWriter)
		} else {
			logrus.Warn("Failed to log to file, using default stderr")
			logrus.SetOutput(os.Stdout)
		}
	} else {
		logrus.SetOutput(os.Stdout)
	}
}
