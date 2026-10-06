package models

import (
	"time"

	"gorm.io/gorm"
)

type Auth struct {
	ID           uint           `gorm:"primaryKey"`
	UserID       uint           `gorm:"uniqueIndex;not null"`
	OIDCSub      *string        `gorm:"column:oidc_sub;size:255;uniqueIndex"`
	PasswordHash string         `gorm:"column:password_hash;size:255"`
	LoginMode    int            `gorm:"column:login_mode;type:int;default:1"`
	LastLoginAt  *time.Time     `gorm:"column:last_login_at"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Auth) TableName() string {
	return "auth"
}
