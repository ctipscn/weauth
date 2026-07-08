package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Challenge 记录表
type Challenge struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	SiteID     uint       `gorm:"not null;index" json:"site_id"`
	Challenge  string     `gorm:"size:128;not null" json:"challenge"` // SHA-256 hex = 64 chars, with margin
	Difficulty int        `gorm:"not null" json:"difficulty"`
	ClientIP   string     `gorm:"size:45" json:"client_ip"`
	Origin     string     `gorm:"size:255" json:"origin"`
	Action     string     `gorm:"size:100" json:"action"`
	Solved     bool       `gorm:"default:false" json:"solved"`
	Token      string     `gorm:"size:128;index" json:"-"` // UUID*2 = 72 chars, with margin
	TokenUsed  bool       `gorm:"default:false" json:"token_used"`
	ExpiresAt  time.Time  `gorm:"not null" json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	SolvedAt   *time.Time `json:"solved_at,omitempty"`

	// 批量模式字段
	BatchID    *uuid.UUID `gorm:"type:uuid;index" json:"batch_id,omitempty"` // 批次ID，同批次的 challenge 共享
	BatchIndex int        `gorm:"default:0" json:"batch_index"`              // 在批次中的序号
	BatchTotal int        `gorm:"default:1" json:"batch_total"`              // 批次总数

	// 关联
	Site Site `gorm:"foreignKey:SiteID" json:"-"`
}

func (Challenge) TableName() string {
	return "challenges"
}

// BeforeCreate 创建前生成 UUID
func (c *Challenge) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// IsExpired 检查是否过期
func (c *Challenge) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// CanSolve 检查是否可以解决
func (c *Challenge) CanSolve() bool {
	return !c.Solved && !c.IsExpired()
}

// MarkSolved 标记为已解决并生成 token
func (c *Challenge) MarkSolved() {
	now := time.Now()
	c.Solved = true
	c.SolvedAt = &now
	c.Token = uuid.New().String() + uuid.New().String()
}

// VerificationLog 验证日志表
type VerificationLog struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ChallengeID *uuid.UUID `gorm:"type:uuid;index" json:"challenge_id"`
	SiteID      *uint      `gorm:"index" json:"site_id"`
	ClientIP    string     `gorm:"size:45" json:"client_ip"`
	Origin      string     `gorm:"size:255" json:"origin"`
	Action      string     `gorm:"size:100" json:"action"`
	Iterations  int64      `json:"iterations"`
	SolveTimeMs int        `json:"solve_time_ms"`
	Success     bool       `gorm:"default:false" json:"success"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (VerificationLog) TableName() string {
	return "verification_logs"
}
