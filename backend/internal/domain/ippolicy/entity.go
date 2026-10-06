package ippolicy

import "time"

// IpPolicy IP 策略实体
type IpPolicy struct {
	ID                  uint      `json:"id"`
	SiteID              *uint     `json:"site_id"`
	TimeWindowSeconds   int       `json:"time_window_seconds"`
	RequestThreshold    int       `json:"request_threshold"`
	DifficultyIncrement int       `json:"difficulty_increment"`
	BatchIncrement      int       `json:"batch_increment"`
	FailureWeight       float64   `json:"failure_weight"`
	Enabled             bool      `json:"enabled"`
	CreatedAt           time.Time `json:"created_at"`
}

// NewDefaultIpPolicy 创建默认全局 IP 策略
func NewDefaultIpPolicy() *IpPolicy {
	return &IpPolicy{
		TimeWindowSeconds:   300,
		RequestThreshold:    10,
		DifficultyIncrement: 1,
		BatchIncrement:      1,
		FailureWeight:       2.0,
		Enabled:             true,
	}
}

// NewSiteIpPolicy 创建站点 IP 策略
func NewSiteIpPolicy(siteID uint) *IpPolicy {
	return &IpPolicy{
		SiteID:              &siteID,
		TimeWindowSeconds:   300,
		RequestThreshold:    10,
		DifficultyIncrement: 1,
		BatchIncrement:      1,
		FailureWeight:       2.0,
		Enabled:             true,
	}
}

// GetEffectiveDifficulty 计算有效难度（动态难度模式）
func (p *IpPolicy) GetEffectiveDifficulty(
	baseDifficulty, 
	maxDifficulty, 
	requestCount, 
	failureCount int,
) int {
	
	if !p.Enabled {
		return baseDifficulty
	}

	weightedCount := float64(requestCount) + float64(failureCount)*p.FailureWeight
	increments := int(weightedCount) / p.RequestThreshold
	effectiveDifficulty := baseDifficulty + increments*p.DifficultyIncrement

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

	weightedCount := float64(requestCount) + float64(failureCount)*p.FailureWeight
	increments := int(weightedCount) / p.RequestThreshold
	effectiveBatchCount := minBatchCount + increments*p.BatchIncrement

	if effectiveBatchCount > maxBatchCount {
		return maxBatchCount
	}
	return effectiveBatchCount
}