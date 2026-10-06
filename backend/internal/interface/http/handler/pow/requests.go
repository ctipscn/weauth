package pow

// CreateChallengeRequest 创建 challenge 请求
type CreateChallengeRequest struct {
	Sitekey    string `json:"sitekey" binding:"required"`
	Origin     string `json:"origin" binding:"required"`
	Action     string `json:"action"`
}

// VerifyChallengeRequest 验证 challenge 请求
type VerifyChallengeRequest struct {
	ChallengeID string   `json:"challenge_id" binding:"required"`
	Nonce       string   `json:"nonce"`
	Nonces      []string `json:"nonces"`
	Iterations  int64    `json:"iterations"`
	SolveTimeMs int      `json:"solve_time_ms"`
	BatchMode   bool     `json:"batch_mode"`
}