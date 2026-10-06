package iprecord

import (
	"github.com/ctipscn/weauth/backend/internal/domain/iprecord"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
)

// ToDomain 将 IP 记录模型转换为领域实体
func ToDomain(m *models.IpRecord) *iprecord.IpRecord {
	return &iprecord.IpRecord{
		ID:                m.ID,
		SiteID:            m.SiteID,
		ClientIP:          m.ClientIP,
		RequestCount:      m.RequestCount,
		FailureCount:      m.FailureCount,
		CurrentDifficulty: m.CurrentDifficulty,
		WindowStart:       m.WindowStart,
		LastRequestAt:     m.LastRequestAt,
		CreatedAt:         m.CreatedAt,
	}
}

// ToModel 将 IP 记录领域实体转换为模型
func ToModel(d *iprecord.IpRecord) *models.IpRecord {
	return &models.IpRecord{
		ID:                d.ID,
		SiteID:            d.SiteID,
		ClientIP:          d.ClientIP,
		RequestCount:      d.RequestCount,
		FailureCount:      d.FailureCount,
		CurrentDifficulty: d.CurrentDifficulty,
		WindowStart:       d.WindowStart,
		LastRequestAt:     d.LastRequestAt,
		CreatedAt:         d.CreatedAt,
	}
}