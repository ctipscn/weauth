package types

import (
	"os"
	"sync"

	"github.com/ctipscn/weauth/backend/internal/config/utils"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host                  string
	Port                  string
	User                  string
	Password              string
	DBName                string
	Mode                  string
	MaxIdleConns          int
	MaxOpenConns          int
	ConnMaxLifetimeMinute int
	ConnMaxIdleMinute     int
}

var (
	dbInstance *DatabaseConfig
	dbOnce     sync.Once
)

// Database 返回全局单例 DatabaseConfig。
// 必须先调用 NewDefaultDatabaseConfig 初始化，否则返回 nil。
func Database() *DatabaseConfig {
	return dbInstance
}

// ApplyFormFile 从 INI 文件的 [database] 节加载配置并合并到当前实例
func (c *DatabaseConfig) ApplyFormFile(path string) (*DatabaseConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return c, nil
	}
	sections, err := utils.ParseINI(path)
	if err != nil {
		return nil, err
	}
	if kv, ok := sections["database"]; ok {
		utils.ApplyString(kv, "host", &c.Host)
		utils.ApplyString(kv, "port", &c.Port)
		utils.ApplyString(kv, "user", &c.User)
		utils.ApplyString(kv, "password", &c.Password)
		utils.ApplyString(kv, "db_name", &c.DBName)
		utils.ApplyString(kv, "db_mode", &c.Mode)
		utils.ApplyInt(kv, "max_idle_conns", &c.MaxIdleConns)
		utils.ApplyInt(kv, "max_open_conns", &c.MaxOpenConns)
		utils.ApplyInt(kv, "conn_max_lifetime_minute", &c.ConnMaxLifetimeMinute)
		utils.ApplyInt(kv, "conn_max_idle_minute", &c.ConnMaxIdleMinute)
	}
	return c, nil
}

// ApplyFromENV 从 WEAUTH_ 前缀的环境变量加载配置并合并到当前实例。
// 支持链式调用。
func (c *DatabaseConfig) ApplyFromENV() *DatabaseConfig {
	utils.EnvString("WEAUTH_DATABASE_HOST", &c.Host)
	utils.EnvString("WEAUTH_DATABASE_PORT", &c.Port)
	utils.EnvString("WEAUTH_DATABASE_USER", &c.User)
	utils.EnvString("WEAUTH_DATABASE_PASSWORD", &c.Password)
	utils.EnvString("WEAUTH_DATABASE_DB_NAME", &c.DBName)
	utils.EnvString("WEAUTH_DATABASE_DB_MODE", &c.Mode)
	utils.EnvInt("WEAUTH_DATABASE_MAX_IDLE_CONNS", &c.MaxIdleConns)
	utils.EnvInt("WEAUTH_DATABASE_MAX_OPEN_CONNS", &c.MaxOpenConns)
	utils.EnvInt("WEAUTH_DATABASE_CONN_MAX_LIFETIME_MINUTE", &c.ConnMaxLifetimeMinute)
	utils.EnvInt("WEAUTH_DATABASE_CONN_MAX_IDLE_MINUTE", &c.ConnMaxIdleMinute)
	return c
}

// NewDefaultDatabaseConfig 初始化并返回全局单例 DatabaseConfig
func NewDefaultDatabaseConfig() *DatabaseConfig {
	dbOnce.Do(func() {
		dbInstance = &DatabaseConfig{
			Host:                  "127.0.0.1",
			Port:                  "5432",
			User:                  "postgres",
			Password:              "admin123",
			Mode:                  "release",
			DBName:                "pow_verify",
			MaxIdleConns:          10,
			MaxOpenConns:          100,
			ConnMaxLifetimeMinute: 30,
			ConnMaxIdleMinute:     5,
		}
	})
	return dbInstance
}

// DSN 返回 PostgreSQL 连接字符串
func (c *DatabaseConfig) DSN() string {
	return "host=" + c.Host +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" port=" + c.Port +
		" sslmode=disable TimeZone=Asia/Shanghai"
}
