package sql

import (
	pkgConfig "discord-chatbot/pkg/config"
	"discord-chatbot/pkg/logger"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
)

var Database = &gorm.DB{}

var log = logger.GetLogger("sql")

func InitSql() error {

	dbPath := pkgConfig.Args.DiscordDatabasePath

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm); err != nil {
			return fmt.Errorf("create dir failed: %v", err)
		}

		if f, err := os.Create(dbPath); err == nil {
			f.Close()
		} else {
			return fmt.Errorf("create db file failed: %v", err)
		}
		log.Infof("Create new %s file", filepath.Base(dbPath))
	}

	db, err := gorm.Open(gormlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect db failed: %v", err)
	}

	Database = db

	if err := Database.AutoMigrate(&User{}); err != nil {
		return fmt.Errorf("migrate failed: %v", err)
	}
	if err := Database.AutoMigrate(&DiscordGambleGamer{}); err != nil {
		return fmt.Errorf("migrate failed: %v", err)
	}
	return nil
}
