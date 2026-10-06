package site

import (
	"github.com/ctipscn/weauth/backend/internal/domain/site"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
)

// ToDomain 将站点模型转换为领域实体
func ToDomain(m *models.Site) *site.Site {
	domains := make([]site.SiteDomain, len(m.Domains))
	for i, d := range m.Domains {
		domains[i] = site.SiteDomain{
			ID:        d.ID,
			SiteID:    d.SiteID,
			Domain:    d.Domain,
			CreatedAt: d.CreatedAt,
		}
	}

	return &site.Site{
		ID:               m.ID,
		UserID:           m.UserID,
		Name:             m.Name,
		Sitekey:          m.Sitekey,
		Secret:           m.Secret,
		BaseDifficulty:   m.BaseDifficulty,
		MaxDifficulty:    m.MaxDifficulty,
		ChallengeTimeout: m.ChallengeTimeout,
		TokenTimeout:     m.TokenTimeout,
		Enabled:          m.Enabled,
		BatchModeEnabled: m.BatchModeEnabled,
		BatchDifficulty:  m.BatchDifficulty,
		MinBatchCount:    m.MinBatchCount,
		MaxBatchCount:    m.MaxBatchCount,
		Domains:          domains,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

// ToModel 将站点领域实体转换为模型
func ToModel(d *site.Site) *models.Site {
	return &models.Site{
		ID:               d.ID,
		UserID:           d.UserID,
		Name:             d.Name,
		Sitekey:          d.Sitekey,
		Secret:           d.Secret,
		BaseDifficulty:   d.BaseDifficulty,
		MaxDifficulty:    d.MaxDifficulty,
		ChallengeTimeout: d.ChallengeTimeout,
		TokenTimeout:     d.TokenTimeout,
		Enabled:          d.Enabled,
		BatchModeEnabled: d.BatchModeEnabled,
		BatchDifficulty:  d.BatchDifficulty,
		MinBatchCount:    d.MinBatchCount,
		MaxBatchCount:    d.MaxBatchCount,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
}