package pow

// CreateChallengeRequest 创建 challenge 请求
type CreateChallengeRequest struct {
	Sitekey    string `json:"sitekey" binding:"required"`
	Origin     string `json:"origin" binding:"required"`
	Action     string `json:"action"`
	Difficulty *int   `json:"difficulty"` // 可选，用于测试指定难度
}

// VerifyChallengeRequest 验证 challenge 请求
type VerifyChallengeRequest struct {
	ChallengeID string   `json:"challenge_id" binding:"required"`
	Nonce       string   `json:"nonce"`   // 单个模式使用
	Nonces      []string `json:"nonces"`  // 批量模式使用
	Iterations  int64    `json:"iterations"`
	SolveTimeMs int      `json:"solve_time_ms"`
	BatchMode   bool     `json:"batch_mode"` // 是否为批量模式
}