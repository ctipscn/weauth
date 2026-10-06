package verifylog

import (
	"github.com/ctipscn/weauth/backend/internal/domain/verifylog"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
)

// ToDomain 将验证日志模型转换为领域实体
func ToDomain(m *models.VerificationLog) *verifylog.VerificationLog {
	return &verifylog.VerificationLog{
		ID:          m.ID,
		ChallengeID: m.ChallengeID,
		SiteID:      m.SiteID,
		ClientIP:    m.ClientIP,
		Origin:      m.Origin,
		Action:      m.Action,
		Iterations:  m.Iterations,
		SolveTimeMs: m.SolveTimeMs,
		Success:     m.Success,
		CreatedAt:   m.CreatedAt,
	}
}

// ToModel 将验证日志领域实体转换为模型
func ToModel(d *verifylog.VerificationLog) *models.VerificationLog {
	return &models.VerificationLog{
		ID:          d.ID,
		ChallengeID: d.ChallengeID,
		SiteID:      d.SiteID,
		ClientIP:    d.ClientIP,
		Origin:      d.Origin,
		Action:      d.Action,
		Iterations:  d.Iterations,
		SolveTimeMs: d.SolveTimeMs,
		Success:     d.Success,
		CreatedAt:   d.CreatedAt,
	}
}