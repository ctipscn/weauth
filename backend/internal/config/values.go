package config

// Load() 使用的环境变量键。
const (
	EnvServerPort       = "SERVER_PORT"
	EnvServerMode       = "SERVER_MODE"
	EnvDBHost           = "DB_HOST"
	EnvDBPort           = "DB_PORT"
	EnvDBUser           = "DB_USER"
	EnvDBPassword       = "DB_PASSWORD"
	EnvDBName           = "DB_NAME"
	EnvDBDriver         = "DB_DRIVER"
	EnvDBPath           = "DB_PATH"
	EnvJWTSecret        = "JWT_SECRET"
	EnvJWTExpireHours   = "JWT_EXPIRE_HOURS"
	EnvOIDCClientID     = "OIDC_CLIENT_ID"
	EnvOIDCClientSecret = "OIDC_CLIENT_SECRET"
	EnvOIDCDiscoveryURL = "OIDC_DISCOVERY_URL"
	EnvOIDCRedirectURL  = "OIDC_REDIRECT_URL"
	EnvFrontendURL      = "FRONTEND_URL"
	EnvCORSOrigins      = "CORS_ORIGINS"
)

// 环境变量为空时 Load() 使用的默认值。
const (
	DefaultServerPort       = "8080"
	DefaultServerMode       = "debug"
	DefaultDBDriver         = "sqlite"
	DefaultDBHost           = "localhost"
	DefaultDBPort           = "5432"
	DefaultDBUser           = "postgres"
	DefaultDBPassword       = ""
	DefaultDBName           = "pow_verify"
	DefaultDBPath           = "data/weauth.db"
	DefaultJWTSecret        = "weauth-jwt-secret-key"
	DefaultJWTExpireHours   = "24"
	DefaultOIDCClientID     = ""
	DefaultOIDCClientSecret = ""
	DefaultOIDCDiscoveryURL = ""
	DefaultOIDCRedirectURL  = ""
	DefaultFrontendURL      = "http://localhost:5173"
	DefaultCORSOrigins      = "http://localhost:5173"
)


const (
	PostgresDBDriver = "postgres"
	SqliteDBDriver   = "sqlite"
)
