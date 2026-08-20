package auth

// LoginRequest 密码登录请求
type LoginRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	Password string `json:"password"`
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}