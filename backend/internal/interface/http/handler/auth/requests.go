package auth

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 邮箱/用户名密码登录请求
type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// UpdatePasswordRequest 更新密码请求（需要原密码）
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// SetPasswordRequest 设置密码请求（OIDC 用户首次设置密码，不需要原密码）
type SetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}