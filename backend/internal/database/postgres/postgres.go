package postgres

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresStorer 基于 PostgreSQL 实现 database.Storer 接口。
type PostgresStorer struct {
	db *gorm.DB
}

// New 打开 PostgreSQL 连接并返回 PostgresStorer。
func New(dsn string, logLevel logger.LogLevel) (*PostgresStorer, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	log.Println("PostgreSQL connected")
	return &PostgresStorer{db: db}, nil
}

func (s *PostgresStorer) DB() *gorm.DB                       { return s.db }
func (s *PostgresStorer) AutoMigrate(models ...any) error     { return s.db.AutoMigrate(models...) }
func (s *PostgresStorer) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
