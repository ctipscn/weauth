package oidc

// OIDCConfig OIDC 发现配置
type OIDCConfig struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
}

// TokenResponse OIDC token 响应
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
}

// OIDCUserInfo OIDC 用户信息
type OIDCUserInfo struct {
	Sub      string `json:"sub"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
}