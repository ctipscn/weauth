package pow

// CreateChallengeResponse 创建 challenge 响应
type CreateChallengeResponse struct {
	ChallengeID string   `json:"challenge_id"`
	Challenge   string   `json:"challenge,omitempty"`
	Challenges  []string `json:"challenges,omitempty"`
	Difficulty  int      `json:"difficulty"`
	BatchMode   bool     `json:"batch_mode"`
	BatchTotal  int      `json:"batch_total,omitempty"`
}

// VerifyChallengeResponse 验证 challenge 响应
type VerifyChallengeResponse struct {
	Token string `json:"token"`
}