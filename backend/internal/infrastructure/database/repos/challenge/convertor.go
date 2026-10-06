package challenge

import (
	"github.com/ctipscn/weauth/backend/internal/domain/challenge"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
)

// ToDomain 将挑战模型转换为领域实体
func ToDomain(m *models.Challenge) *challenge.Challenge {
	return &challenge.Challenge{
		ID:         m.ID,
		SiteID:     m.SiteID,
		Challenge:  m.Challenge,
		Difficulty: m.Difficulty,
		ClientIP:   m.ClientIP,
		Origin:     m.Origin,
		Action:     m.Action,
		Solved:     m.Solved,
		Token:      m.Token,
		TokenUsed:  m.TokenUsed,
		ExpiresAt:  m.ExpiresAt,
		CreatedAt:  m.CreatedAt,
		SolvedAt:   m.SolvedAt,
		BatchID:    m.BatchID,
		BatchIndex: m.BatchIndex,
		BatchTotal: m.BatchTotal,
	}
}

// ToModel 将挑战领域实体转换为模型
func ToModel(d *challenge.Challenge) *models.Challenge {
	return &models.Challenge{
		ID:         d.ID,
		SiteID:     d.SiteID,
		Challenge:  d.Challenge,
		Difficulty: d.Difficulty,
		ClientIP:   d.ClientIP,
		Origin:     d.Origin,
		Action:     d.Action,
		Solved:     d.Solved,
		Token:      d.Token,
		TokenUsed:  d.TokenUsed,
		ExpiresAt:  d.ExpiresAt,
		CreatedAt:  d.CreatedAt,
		SolvedAt:   d.SolvedAt,
		BatchID:    d.BatchID,
		BatchIndex: d.BatchIndex,
		BatchTotal: d.BatchTotal,
	}
}
