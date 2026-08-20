package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)



// Config 顶层应用配置结构体。
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OIDC     OIDCConfig
	CORS     CORSConfig
}

// AppConfig 应用配置单例。
var AppConfig *Config

// Load 从环境变量（以及 .env 文件，如果存在）读取配置。
func Load() error {
	_ = godotenv.Load()

	expireHours, _ := strconv.Atoi(getEnv(EnvJWTExpireHours, DefaultJWTExpireHours))

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv(EnvServerPort, DefaultServerPort),
			Mode: getEnv(EnvServerMode, DefaultServerMode),
		},
		Database: DatabaseConfig{
			Driver:   getEnv(EnvDBDriver, DefaultDBDriver),
			Host:     getEnv(EnvDBHost, DefaultDBHost),
			Port:     getEnv(EnvDBPort, DefaultDBPort),
			User:     getEnv(EnvDBUser, DefaultDBUser),
			Password: getEnv(EnvDBPassword, DefaultDBPassword),
			DBName:   getEnv(EnvDBName, DefaultDBName),
			Path:     getEnv(EnvDBPath, DefaultDBPath),
		},
		JWT: JWTConfig{
			Secret:      getEnv(EnvJWTSecret, DefaultJWTSecret),
			ExpireHours: expireHours,
		},
		OIDC: OIDCConfig{
			ClientID:     getEnv(EnvOIDCClientID, DefaultOIDCClientID),
			ClientSecret: getEnv(EnvOIDCClientSecret, DefaultOIDCClientSecret),
			DiscoveryURL: getEnv(EnvOIDCDiscoveryURL, DefaultOIDCDiscoveryURL),
			RedirectURL:  getEnv(EnvOIDCRedirectURL, DefaultOIDCRedirectURL),
			FrontendURL:  getEnv(EnvFrontendURL, DefaultFrontendURL),
		},
		CORS: CORSConfig{
			AllowOrigins: strings.Split(getEnv(EnvCORSOrigins, DefaultCORSOrigins), ","),
		},
	}

	return nil
}

// getEnv 读取环境变量，未设置时返回默认值。
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
