package pow

// CreateChallengeResponse 创建 challenge 响应
type CreateChallengeResponse struct {
	ChallengeID string   `json:"challenge_id"`
	Challenge   string   `json:"challenge,omitempty"`   // 单个模式使用
	Challenges  []string `json:"challenges,omitempty"`  // 批量模式使用
	Difficulty  int      `json:"difficulty"`
	BatchMode   bool     `json:"batch_mode"`            // 是否为批量模式
	BatchTotal  int      `json:"batch_total,omitempty"` // 批量模式的总数量
}

// VerifyChallengeResponse 验证 challenge 响应
type VerifyChallengeResponse struct {
	Token string `json:"token"`
}