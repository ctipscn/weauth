package challenge

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Challenge PoW 挑战实体
type Challenge struct {
	ID         uuid.UUID
	SiteID     uint
	Challenge  string
	Difficulty int
	ClientIP   string
	Origin     string
	Action     string
	Solved     bool
	Token      string
	TokenUsed  bool
	ExpiresAt  time.Time
	CreatedAt  time.Time
	SolvedAt   *time.Time

	BatchID    *uuid.UUID
	BatchIndex int
	BatchTotal int
}

// NewChallenge 创建新挑战
func NewChallenge(
	siteID uint, challenge string,
	difficulty int, clientIP, origin,
	action string, expiresAt time.Time) *Challenge {

	return &Challenge{
		ID:         uuid.New(),
		SiteID:     siteID,
		Challenge:  challenge,
		Difficulty: difficulty,
		ClientIP:   clientIP,
		Origin:     origin,
		Action:     action,
		ExpiresAt:  expiresAt,
		BatchTotal: 1,
	}
}

// IsExpired 检查是否过期
func (c *Challenge) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// CanSolve 检查是否可以解决
func (c *Challenge) CanSolve() bool {
	return !c.Solved && !c.IsExpired()
}

// MarkSolved 标记为已解决并生成 token
func (c *Challenge) MarkSolved() {
	now := time.Now()
	c.Solved = true
	c.SolvedAt = &now
	c.Token = uuid.New().String() + uuid.New().String()
}

// VerifySolution 验证解决方案
func (c *Challenge) VerifySolution(nonce string, difficulty int) bool {
	data := c.Challenge + nonce
	hash := sha256.Sum256([]byte(data))
	hashHex := hex.EncodeToString(hash[:])
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hashHex, prefix)
}
