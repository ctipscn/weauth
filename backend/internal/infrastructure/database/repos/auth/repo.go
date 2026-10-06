package auth

import (
	"errors"

	"github.com/ctipscn/weauth/backend/internal/infrastructure/database"
	"gorm.io/gorm"

	"github.com/ctipscn/weauth/backend/internal/domain"
	basedomain "github.com/ctipscn/weauth/backend/internal/domain/auth"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
)

type AuthCredentialRepo struct {
	db *gorm.DB
}

func NewAuthCredentialRepo(driver *database.PostgresDriver) *AuthCredentialRepo {
	return &AuthCredentialRepo{db: driver.DB()}
}

func (r *AuthCredentialRepo) FindByUserID(userID uint) (*basedomain.Auth, error) {
	var m models.Auth
	if err := r.db.Where("user_id = ?", userID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

func (r *AuthCredentialRepo) FindByUsername(username string) (*basedomain.Auth, error) {
	var m models.Auth
	if err := r.db.Joins("JOIN users ON auth.user_id = users.id").
		Where("users.username = ?", username).
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

func (r *AuthCredentialRepo) FindByOIDCSub(sub string) (*basedomain.Auth, error) {
	var m models.Auth
	if err := r.db.Where("oidc_sub = ?", sub).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

func (r *AuthCredentialRepo) Create(credential *basedomain.Auth) error {
	m := ToModel(credential)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	credential.ID = m.ID
	return nil
}

func (r *AuthCredentialRepo) Save(credential *basedomain.Auth) error {
	m := ToModel(credential)
	return r.db.Save(m).Error
}

func (r *AuthCredentialRepo) DeleteByUserID(userID uint) error {
	result := r.db.Where("user_id = ?", userID).Delete(&models.Auth{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
