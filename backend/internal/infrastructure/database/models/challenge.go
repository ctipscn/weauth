package models

import (
	"time"

	"github.com/google/uuid"
)

// Challenge 记录表
type Challenge struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	SiteID     uint       `gorm:"not null;index"`
	Challenge  string     `gorm:"size:128;not null"` // SHA-256 hex = 64 chars, with margin
	Difficulty int        `gorm:"not null"`
	ClientIP   string     `gorm:"size:45"`
	Origin     string     `gorm:"size:255"`
	Action     string     `gorm:"size:100"`
	Solved     bool       `gorm:"default:false"`
	Token      string     `gorm:"size:128;index"`
	TokenUsed  bool       `gorm:"default:false"`
	ExpiresAt  time.Time  `gorm:"not null"`
	CreatedAt  time.Time
	SolvedAt   *time.Time

	// 批量模式字段
	BatchID    *uuid.UUID `gorm:"type:uuid;index"` // 批次ID，同批次的 challenge 共享
	BatchIndex int        `gorm:"default:0"`        // 在批次中的序号
	BatchTotal int        `gorm:"default:1"`        // 批次总数
}

func (Challenge) TableName() string {
	return "challenges"
}

