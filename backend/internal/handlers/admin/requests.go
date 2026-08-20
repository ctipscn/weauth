package admin

// CreateSiteRequest 创建站点请求
type CreateSiteRequest struct {
	Name             string   `json:"name" binding:"required"`
	Domains          []string `json:"domains"`
	BatchModeEnabled *bool    `json:"batch_mode_enabled"` // 默认 true（批量模式）
	BatchDifficulty  *int     `json:"batch_difficulty"`   // 批量模式难度，默认 4
	MinBatchCount    *int     `json:"min_batch_count"`    // 最小批次数，默认 1
	MaxBatchCount    *int     `json:"max_batch_count"`    // 最大批次数，默认 10
}

// UpdateSiteRequest 更新站点请求
type UpdateSiteRequest struct {
	Name             *string `json:"name"`
	BaseDifficulty   *int    `json:"base_difficulty"`
	MaxDifficulty    *int    `json:"max_difficulty"`
	ChallengeTimeout *int    `json:"challenge_timeout"`
	TokenTimeout     *int    `json:"token_timeout"`
	Enabled          *bool   `json:"enabled"`
	// 批量模式配置
	BatchModeEnabled *bool `json:"batch_mode_enabled"`
	BatchDifficulty  *int  `json:"batch_difficulty"`
	MinBatchCount    *int  `json:"min_batch_count"`
	MaxBatchCount    *int  `json:"max_batch_count"`
}

// AddDomainRequest 添加域名请求
type AddDomainRequest struct {
	Domain string `json:"domain" binding:"required"`
}

// UpdateIpPolicyRequest 更新 IP 策略请求
type UpdateIpPolicyRequest struct {
	TimeWindowSeconds   *int     `json:"time_window_seconds"`
	RequestThreshold    *int     `json:"request_threshold"`
	DifficultyIncrement *int     `json:"difficulty_increment"` // 动态难度模式：难度增量
	BatchIncrement      *int     `json:"batch_increment"`      // 批量模式：批次增量
	FailureWeight       *float64 `json:"failure_weight"`
	Enabled             *bool    `json:"enabled"`
}