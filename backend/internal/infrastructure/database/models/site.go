package models

import (
	"time"

	"gorm.io/gorm"
)

// Site 站点表
type Site struct {
	ID               uint   `gorm:"primaryKey"`
	UserID           uint   `gorm:"index;default:0"`
	Name             string `gorm:"size:100;not null"`
	Sitekey          string `gorm:"size:32;uniqueIndex;not null"`
	Secret           string `gorm:"size:64;uniqueIndex;not null"`
	BaseDifficulty   int    `gorm:"default:4;not null"`
	MaxDifficulty    int    `gorm:"default:8;not null"`
	ChallengeTimeout int    `gorm:"default:120"`
	TokenTimeout     int    `gorm:"default:60"` // 默认 1 分钟
	Enabled          bool   `gorm:"default:true"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	// 批量模式配置
	BatchModeEnabled bool `gorm:"default:false"` // 是否启用批量模式
	BatchDifficulty  int  `gorm:"default:4"`     // 批量模式固定难度
	MinBatchCount    int  `gorm:"default:1"`     // 最小批次数
	MaxBatchCount    int  `gorm:"default:10"`    // 最大批次数

	// 关联域名
	Domains []Domain `gorm:"foreignKey:SiteID"`
}

// TableName 设置表名
func (Site) TableName() string {
	return "sites"
}