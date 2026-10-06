package auth

import (
	domain "github.com/ctipscn/weauth/backend/internal/domain/auth"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
)

func ToDomain(m *models.Auth) *domain.Auth {
	if m == nil {
		return nil
	}
	return &domain.Auth{
		ID:           m.ID,
		UserID:       m.UserID,
		OIDCSub:      m.OIDCSub,
		PasswordHash: m.PasswordHash,
		LoginMode:    domain.LoginMode(m.LoginMode),
		LastLoginAt:  m.LastLoginAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func ToModel(d *domain.Auth) *models.Auth {
	if d == nil {
		return nil
	}
	return &models.Auth{
		ID:           d.ID,
		UserID:       d.UserID,
		OIDCSub:      d.OIDCSub,
		PasswordHash: d.PasswordHash,
		LoginMode:    int(d.LoginMode),
		LastLoginAt:  d.LastLoginAt,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}
