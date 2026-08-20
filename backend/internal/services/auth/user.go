package auth

// UserInfo OIDC 用户信息
type UserInfo struct {
	Sub      string `json:"sub"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
}