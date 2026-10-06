package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID        uint           `gorm:"primaryKey"`
	Username  string         `gorm:"size:100"`
	Email     string         `gorm:"size:255"`
	Avatar    string         `gorm:"size:500"`
	Role      string         `gorm:"size:20;default:staff"`
	Enabled   bool           `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}