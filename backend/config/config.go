package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OIDC     OIDCConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type OIDCConfig struct {
	ClientID     string
	ClientSecret string
	DiscoveryURL string
	RedirectURL  string
	FrontendURL  string
}

type CORSConfig struct {
	AllowOrigins []string
}

var AppConfig *Config

func Load() error {
	// 尝试加载 .env 文件，如果不存在则使用环境变量
	_ = godotenv.Load()

	expireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "pow_verify"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "default-secret-please-change"),
			ExpireHours: expireHours,
		},
	OIDC: OIDCConfig{
		ClientID:     getEnv("OIDC_CLIENT_ID", ""),
		ClientSecret: getEnv("OIDC_CLIENT_SECRET", ""),
		DiscoveryURL: getEnv("OIDC_DISCOVERY_URL", ""),
		RedirectURL:  getEnv("OIDC_REDIRECT_URL", ""),
		FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:5173"),
	},
	CORS: CORSConfig{
		// 开发环境：http://localhost:5173
		// 生产环境：https://www.weauth.cn,https://api.weauth.cn
		// 允许 www.weauth.cn 前端访问，以及 iframe widget 页面调用 API
		AllowOrigins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173"), ","),
	},
}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *DatabaseConfig) DSN() string {
	return "host=" + c.Host +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" port=" + c.Port +
		" sslmode=disable TimeZone=Asia/Shanghai"
}
