package ippolicy

import (
	"github.com/ctipscn/weauth/backend/internal/domain/ippolicy"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
)

// ToDomain 将 IP 策略模型转换为领域实体
func ToDomain(m *models.IpPolicy) *ippolicy.IpPolicy {
	return &ippolicy.IpPolicy{
		ID:                  m.ID,
		SiteID:              m.SiteID,
		TimeWindowSeconds:   m.TimeWindowSeconds,
		RequestThreshold:    m.RequestThreshold,
		DifficultyIncrement: m.DifficultyIncrement,
		BatchIncrement:      m.BatchIncrement,
		FailureWeight:       m.FailureWeight,
		Enabled:             m.Enabled,
		CreatedAt:           m.CreatedAt,
	}
}

// ToModel 将 IP 策略领域实体转换为模型
func ToModel(d *ippolicy.IpPolicy) *models.IpPolicy {
	return &models.IpPolicy{
		ID:                  d.ID,
		SiteID:              d.SiteID,
		TimeWindowSeconds:   d.TimeWindowSeconds,
		RequestThreshold:    d.RequestThreshold,
		DifficultyIncrement: d.DifficultyIncrement,
		BatchIncrement:      d.BatchIncrement,
		FailureWeight:       d.FailureWeight,
		Enabled:             d.Enabled,
		CreatedAt:           d.CreatedAt,
	}
}