package cli

import (
	"github.com/ctipscn/weauth/backend/internal/config"
)

// ConfigOfArgs 保存命令行 --xx-x 参数指定的配置覆盖值。
type ConfigOfArgs struct {
	ServerListen               *string
	ServerPort                 *string
	ServerMode                 *string
	ServerLogFile              *string
	DatabaseHost               *string
	DatabasePort               *string
	DatabaseMode               *string
	DatabaseUser               *string
	DatabasePassword           *string
	DatabaseDBName             *string
	DatabaseMaxIdleConns       *int
	DatabaseMaxOpenConns       *int
	DatabaseConnMaxLifetimeMin *int
	DatabaseConnMaxIdleMin     *int
	JWTSecret                  *string
	JWTExpireHours             *int
	OIDCClientID               *string
	OIDCClientSecret           *string
	OIDCDiscoveryURL           *string
	OIDCRedirectURL            *string
	OIDCFrontendURL            *string
	CORSAllowOrigins           *string

	// RateLimit
	RateLimitTimeWindowSeconds *int
	RateLimitRequestLimit      *int
}

// ApplyFromArgs 将命令行覆盖值写入 AppConfig，覆盖优先级高于 INI 和环境变量。
func (f ConfigOfArgs) ApplyFromArgs(cfg *config.AppConfig) {
	// Server
	assignString(f.ServerListen, &cfg.Server.Listen)
	assignString(f.ServerPort, &cfg.Server.Port)
	assignString(f.ServerMode, &cfg.Server.Mode)
	assignString(f.ServerLogFile, &cfg.Server.LogFile)
	// Database
	assignString(f.DatabaseHost, &cfg.Database.Host)
	assignString(f.DatabasePort, &cfg.Database.Port)
	assignString(f.DatabaseMode, &cfg.Database.Mode)
	assignString(f.DatabaseUser, &cfg.Database.User)
	assignString(f.DatabasePassword, &cfg.Database.Password)
	assignString(f.DatabaseDBName, &cfg.Database.DBName)
	assignInt(f.DatabaseMaxIdleConns, &cfg.Database.MaxIdleConns)
	assignInt(f.DatabaseMaxOpenConns, &cfg.Database.MaxOpenConns)
	assignInt(f.DatabaseConnMaxLifetimeMin, &cfg.Database.ConnMaxLifetimeMinute)
	assignInt(f.DatabaseConnMaxIdleMin, &cfg.Database.ConnMaxIdleMinute)
	// JWT
	assignString(f.JWTSecret, &cfg.JWT.Secret)
	assignInt(f.JWTExpireHours, &cfg.JWT.ExpireHours)
	// OIDC
	assignString(f.OIDCClientID, &cfg.OIDC.ClientID)
	assignString(f.OIDCClientSecret, &cfg.OIDC.ClientSecret)
	assignString(f.OIDCDiscoveryURL, &cfg.OIDC.DiscoveryURL)
	assignString(f.OIDCRedirectURL, &cfg.OIDC.RedirectURL)
	assignString(f.OIDCFrontendURL, &cfg.OIDC.FrontendURL)
	// CORS
	if f.CORSAllowOrigins != nil {
		cfg.CORS.AllowOrigins = parseOriginList(*f.CORSAllowOrigins)
	}
	// RateLimit
	assignInt(f.RateLimitTimeWindowSeconds, &cfg.RateLimit.TimeWindowSeconds)
	assignInt(f.RateLimitRequestLimit, &cfg.RateLimit.RequestLimit)
}