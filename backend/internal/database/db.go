package database

import (
	"log"

	"github.com/ctipscn/weauth/backend/internal/config"
	"github.com/ctipscn/weauth/backend/internal/database/postgres"
	"github.com/ctipscn/weauth/backend/internal/database/sqlite"

	"gorm.io/gorm/logger"
)

// dbStoreInstance 全局数据库存储实例。
var dbStoreInstance Storer

// NewDefault 根据配置打开数据库连接，测试连通性，返回 *Storer 及错误。
// 调用方必须将非 nil 错误视为致命错误并退出。
func NewDefault() (*Storer, error) {
	logLevel := logger.Info
	if config.AppConfig.Server.Mode == "release" {
		logLevel = logger.Warn
	}

	cfg := config.AppConfig.Database

	var err error

	if cfg.Driver == "sqlite" {
		dbStoreInstance, err = sqlite.New(cfg.Path, logLevel)
	} else {
		dbStoreInstance, err = postgres.New(cfg.DSN(), logLevel)
	}

	if err != nil {
		return nil, err
	}

	// 测试数据库连接
	db := dbStoreInstance.DB()
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return &dbStoreInstance, nil
}
