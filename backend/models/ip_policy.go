package models

import (
	"time"
)

// IpPolicy IP 策略配置表
type IpPolicy struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	SiteID              *uint     `gorm:"uniqueIndex" json:"site_id"` // NULL 表示全局策略
	TimeWindowSeconds   int       `gorm:"default:300" json:"time_window_seconds"`
	RequestThreshold    int       `gorm:"default:10" json:"request_threshold"`
	DifficultyIncrement int       `gorm:"default:1" json:"difficulty_increment"` // 动态难度模式：每次增加的难度
	BatchIncrement      int       `gorm:"default:1" json:"batch_increment"`      // 批量模式：每次增加的批次数
	FailureWeight       float64   `gorm:"type:decimal(3,2);default:2.0" json:"failure_weight"`
	Enabled             bool      `gorm:"default:true" json:"enabled"`
	CreatedAt           time.Time `json:"created_at"`

	// 关联
	Site *Site `gorm:"foreignKey:SiteID" json:"-"`
}

func (IpPolicy) TableName() string {
	return "ip_policies"
}

// GetEffectiveDifficulty 计算有效难度（动态难度模式）
func (p *IpPolicy) GetEffectiveDifficulty(baseDifficulty, maxDifficulty, requestCount, failureCount int) int {
	if !p.Enabled {
		return baseDifficulty
	}

	// 计算加权请求数
	weightedCount := float64(requestCount) + float64(failureCount)*p.FailureWeight

	// 计算难度增量
	increments := int(weightedCount) / p.RequestThreshold
	effectiveDifficulty := baseDifficulty + increments*p.DifficultyIncrement

	// 不超过最大难度
	if effectiveDifficulty > maxDifficulty {
		return maxDifficulty
	}

	return effectiveDifficulty
}

// GetEffectiveBatchCount 计算有效批次数量（批量模式）
func (p *IpPolicy) GetEffectiveBatchCount(minBatchCount, maxBatchCount, requestCount, failureCount int) int {
	if !p.Enabled {
		return minBatchCount
	}

	// 计算加权请求数
	weightedCount := float64(requestCount) + float64(failureCount)*p.FailureWeight

	// 计算批次增量
	increments := int(weightedCount) / p.RequestThreshold
	effectiveBatchCount := minBatchCount + increments*p.BatchIncrement

	// 不超过最大批次数
	if effectiveBatchCount > maxBatchCount {
		return maxBatchCount
	}

	return effectiveBatchCount
}
