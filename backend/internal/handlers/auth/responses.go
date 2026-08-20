package auth

// LoginResponse 登录响应
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}