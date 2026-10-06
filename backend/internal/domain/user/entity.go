package user

import (
	"time"
)

// User 用户实体
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Role      Role      `json:"role"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsAdmin 检查是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == Administrator
}

// Deactivate 将用户标记为停用状态
func (u *User) Deactivate() {
	u.Enabled = false
}
