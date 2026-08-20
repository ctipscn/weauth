package sqlite

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SqliteStorer 基于 SQLite 实现 database.Storer 接口。
type SqliteStorer struct {
	db *gorm.DB
}

// New 打开 SQLite 数据库文件，必要时创建父目录。
func New(path string, logLevel logger.LogLevel) (*SqliteStorer, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir %s: %w", dir, err)
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	log.Printf("SQLite opened: %s", path)
	return &SqliteStorer{db: db}, nil
}

func (s *SqliteStorer) DB() *gorm.DB                       { return s.db }
func (s *SqliteStorer) AutoMigrate(models ...any) error     { return s.db.AutoMigrate(models...) }
func (s *SqliteStorer) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
