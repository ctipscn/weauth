package database

import (
	"github.com/ctipscn/weauth/backend/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	var err error

	logLevel := logger.Info
	if config.AppConfig.Server.Mode == "release" {
		logLevel = logger.Warn
	}

	DB, err = gorm.Open(postgres.Open(config.AppConfig.Database.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return nil
}

func AutoMigrate(models ...interface{}) error {
	return DB.AutoMigrate(models...)
}
