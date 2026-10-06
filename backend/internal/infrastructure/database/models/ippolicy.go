package models

import (
	"time"
)

// IpPolicy IP 策略配置表
type IpPolicy struct {
	ID                  uint    `gorm:"primaryKey"`
	SiteID              *uint   `gorm:"uniqueIndex"` // NULL 表示全局策略
	TimeWindowSeconds   int     `gorm:"default:300"`
	RequestThreshold    int     `gorm:"default:10"`
	DifficultyIncrement int     `gorm:"default:1"` // 动态难度模式：每次增加的难度
	BatchIncrement      int     `gorm:"default:1"` // 批量模式：每次增加的批次数
	FailureWeight       float64 `gorm:"type:decimal(3,2);default:2.0"`
	Enabled             bool    `gorm:"default:true"`
	CreatedAt           time.Time

}

func (IpPolicy) TableName() string {
	return "ip_policies"
}