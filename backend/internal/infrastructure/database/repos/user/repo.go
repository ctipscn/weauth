package user

import (
	"errors"

	"github.com/ctipscn/weauth/backend/internal/domain"
	"github.com/ctipscn/weauth/backend/internal/domain/user"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

// UserRepo 用户仓储实现
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo 创建用户仓储
func NewUserRepo(driver *database.PostgresDriver) *UserRepo {
	return &UserRepo{db: driver.DB()}
}

// FindByID 根据ID查找用户
func (r *UserRepo) FindByID(id uint) (*user.User, error) {
	var m models.User
	if err := r.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepo) FindByUsername(username string) (*user.User, error) {
	var m models.User
	if err := r.db.Where("username = ?", username).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepo) FindByEmail(email string) (*user.User, error) {
	var m models.User
	if err := r.db.Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// Create 创建用户
func (r *UserRepo) Create(u *user.User) error {
	m := ToModel(u)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	// 回写自增ID
	u.ID = m.ID
	return nil
}

// CountAll 返回用户总数
func (r *UserRepo) CountAll() (int64, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindAll 返回所有用户列表
func (r *UserRepo) FindAll() ([]*user.User, error) {
	var ms []models.User
	if err := r.db.Order("created_at DESC").Find(&ms).Error; err != nil {
		return nil, err
	}
	result := make([]*user.User, 0, len(ms))
	for i := range ms {
		result = append(result, ToDomain(&ms[i]))
	}
	return result, nil
}

// Save 保存用户（全量更新）
func (r *UserRepo) Save(u *user.User) error {
	m := ToModel(u)
	return r.db.Save(m).Error
}

// Delete 删除用户
func (r *UserRepo) Delete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

// Deactivate 软删除用户：设置 enabled=false 并填充 deleted_at
func (r *UserRepo) Deactivate(id uint) error {
	result := r.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"enabled":    false,
		"deleted_at": gorm.Expr("NOW()"),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

// Ensure compile-time interface compliance
var _ user.Repository = (*UserRepo)(nil)
