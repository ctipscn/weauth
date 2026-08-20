package config

// OIDCConfig OIDC 认证配置。
type OIDCConfig struct {
	ClientID     string
	ClientSecret string
	DiscoveryURL string
	RedirectURL  string
	FrontendURL  string
}
