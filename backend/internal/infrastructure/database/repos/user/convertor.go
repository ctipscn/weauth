package user

import (
	"github.com/ctipscn/weauth/backend/internal/domain/user"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
)

// ToDomain 将模型转换为领域实体
func ToDomain(m *models.User) *user.User {
	role, _ := user.ParseRole(m.Role)
	return &user.User{
		ID:        m.ID,
		Username:  m.Username,
		Email:     m.Email,
		Avatar:    m.Avatar,
		Role:      role,
		Enabled:   m.Enabled,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// ToModel 将领域实体转换为模型
func ToModel(d *user.User) *models.User {
	return &models.User{
		ID:        d.ID,
		Username:  d.Username,
		Email:     d.Email,
		Avatar:    d.Avatar,
		Role:      d.Role.String(),
		Enabled:   d.Enabled,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}