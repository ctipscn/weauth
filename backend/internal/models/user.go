package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 管理员用户表
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	OIDCSub      string         `gorm:"column:oidc_sub;size:255;uniqueIndex" json:"oidc_sub,omitempty"`
	Username     string         `gorm:"size:100" json:"username"`
	Email        string         `gorm:"size:255" json:"email"`
	Avatar       string         `gorm:"size:500" json:"avatar"`
	PasswordHash string         `gorm:"size:255" json:"-"`
	Role         string         `gorm:"size:20;default:user" json:"role"`
	Enabled      bool           `gorm:"default:true" json:"enabled"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "admin_users"
}

// SetPassword 设置密码
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// HasPassword 检查是否设置了密码
func (u *User) HasPassword() bool {
	return u.PasswordHash != ""
}

// IsAdmin 检查是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// UpdateLastLogin 更新最后登录时间
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}
