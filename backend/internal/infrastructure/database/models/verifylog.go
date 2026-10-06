package models

import (
	"time"

	"github.com/google/uuid"
)

// VerificationLog 验证日志表
type VerificationLog struct {
	ID          uint       `gorm:"primaryKey"`
	ChallengeID *uuid.UUID `gorm:"type:uuid;index"`
	SiteID      *uint      `gorm:"index"`
	ClientIP    string     `gorm:"size:45"`
	Origin      string     `gorm:"size:255"`
	Action      string     `gorm:"size:100"`
	Iterations  int64
	SolveTimeMs int
	Success     bool `gorm:"default:false"`
	CreatedAt   time.Time
}

func (VerificationLog) TableName() string {
	return "verification_logs"
}
