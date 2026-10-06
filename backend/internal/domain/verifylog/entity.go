package verifylog

import (
	"time"

	"github.com/ctipscn/weauth/backend/internal/domain/challenge"
	"github.com/google/uuid"
)

// VerificationLog 验证日志实体
type VerificationLog struct {
	ID          uint
	ChallengeID *uuid.UUID
	SiteID      *uint
	ClientIP    string
	Origin      string
	Action      string
	Iterations  int64
	SolveTimeMs int
	Success     bool
	CreatedAt   time.Time
}

// NewVerificationLog 创建验证日志
func NewVerificationLog(
	c *challenge.Challenge,
	clientIP string,
	iterations int64,
	solveTimeMs int,
	success bool,
) *VerificationLog {

	return &VerificationLog{
		ChallengeID: &c.ID,
		SiteID:      &c.SiteID,
		ClientIP:    clientIP,
		Origin:      c.Origin,
		Action:      c.Action,
		Iterations:  iterations,
		SolveTimeMs: solveTimeMs,
		Success:     success,
	}
}