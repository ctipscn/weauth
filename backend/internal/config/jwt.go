package config

// JWTConfig JWT 令牌配置。
type JWTConfig struct {
	Secret      string
	ExpireHours int
}
