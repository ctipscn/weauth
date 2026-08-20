package siteverify

// SiteVerifyResponse 站点验证响应
type SiteVerifyResponse struct {
	Success     bool     `json:"success"`
	Hostname    string   `json:"hostname,omitempty"`
	Action      string   `json:"action,omitempty"`
	ChallengeTs string   `json:"challenge_ts,omitempty"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
}