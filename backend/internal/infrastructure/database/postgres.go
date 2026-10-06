package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/ctipscn/weauth/backend/internal/config"
	"github.com/ctipscn/weauth/backend/internal/infrastructure"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// PostgresDriver 封装 *gorm.DB，拒绝直接暴露裸指针
type PostgresDriver struct {
	db  *gorm.DB
	log *infralogger.LogService
}

var (
	driver *PostgresDriver
	once   sync.Once
)

// NewPostgresDriver 初始化并返回单例 PostgresDriver，完成后执行 Ping 并打印结果
func NewPostgresDriver(log *infralogger.LogService) (*PostgresDriver, error) {
	var initErr error

	once.Do(func() {
		cfg := config.Config().Database

		log.Info("[DB] Connecting to database...")

		db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
			Logger: gormlogger.Default.LogMode(gormlogger.Info),
		})
		if err != nil {
			initErr = fmt.Errorf("[DB] Failed to connect database: %w", err)
			log.Error(initErr.Error())
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			initErr = fmt.Errorf("[DB] Failed to get database instance: %w", err)
			log.Error(initErr.Error())
			return
		}

		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeMinute) * time.Minute)
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleMinute) * time.Minute)

		driver = &PostgresDriver{db: db, log: log}

		// 执行 Ping
		if err := driver.Ping(); err != nil {
			initErr = fmt.Errorf("[DB] Database Ping failed: %w", err)
			log.Error(initErr.Error())
			driver = nil
			return
		}

		msg := fmt.Sprintf("[DB] Database connected (host=%s, dbname=%s)", cfg.Host, cfg.DBName)
		// 打印到终端
		fmt.Printf("[PostgresDriver] %s\n", msg)
		// 记录到日志文件
		log.Info(msg)
	})

	if initErr != nil {
		return nil, initErr
	}
	return driver, nil
}

// Ping 检查数据库连接是否正常
func (d *PostgresDriver) Ping() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// DB 返回内部的 *gorm.DB 实例（供仓储层使用）
func (d *PostgresDriver) DB() *gorm.DB {
	return d.db
}

// AutoMigrate 自动迁移数据库表
func (d *PostgresDriver) AutoMigrate(models ...any) error {
	return d.db.AutoMigrate(models...)
}

// compile-time check: PostgresDriver 实现了 Configurator 接口
var _ infrastructure.Configurator = (*PostgresDriver)(nil)

// ApplyMode 根据运行模式切换 GORM 日志级别，不泄漏 GORM 实现细节。
func (d *PostgresDriver) ApplyMode(m infrastructure.Mode) error {
	if d.db == nil {
		return nil
	}

	var level gormlogger.LogLevel
	switch m {
	case infrastructure.ReleaseMode:
		// 生产：仅记录警告和慢查询
		level = gormlogger.Warn
	case infrastructure.TestMode:
		// 测试：仅记录错误
		level = gormlogger.Error
	default:
		// 开发：打印全量 SQL
		level = gormlogger.Info
	}

	d.db.Logger = gormlogger.Default.LogMode(level)
	d.log.Info(fmt.Sprintf("[DB] GORM log level set to %s (mode=%s)",
		gormLogLevelToString(level), m))

	return nil
}

func gormLogLevelToString(level gormlogger.LogLevel) string {
	switch level {
	case gormlogger.Silent:
		return "silent"
	case gormlogger.Error:
		return "error"
	case gormlogger.Warn:
		return "warn"
	case gormlogger.Info:
		return "info"
	default:
		return "unknown"
	}
}
