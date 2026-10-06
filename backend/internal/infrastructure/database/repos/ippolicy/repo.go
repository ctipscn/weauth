package ippolicy

import (
	"errors"

	"github.com/ctipscn/weauth/backend/internal/domain"
	"github.com/ctipscn/weauth/backend/internal/domain/ippolicy"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
	"gorm.io/gorm"
)

// IpPolicyRepo IP 策略仓储实现
type IpPolicyRepo struct {
	db *gorm.DB
}

// NewIpPolicyRepo 创建 IP 策略仓储
func NewIpPolicyRepo(driver *database.PostgresDriver) *IpPolicyRepo {
	return &IpPolicyRepo{db: driver.DB()}
}

// FindBySiteID 查找站点的 IP 策略
func (r *IpPolicyRepo) FindBySiteID(siteID uint) (*ippolicy.IpPolicy, error) {
	var m models.IpPolicy
	if err := r.db.Where("site_id = ?", siteID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// FindGlobal 查找全局 IP 策略
func (r *IpPolicyRepo) FindGlobal() (*ippolicy.IpPolicy, error) {
	var m models.IpPolicy
	if err := r.db.Where("site_id IS NULL").First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// CreateOrUpdate 创建或更新站点的 IP 策略
func (r *IpPolicyRepo) CreateOrUpdate(siteID uint, policy *ippolicy.IpPolicy) (*ippolicy.IpPolicy, error) {
	var existing models.IpPolicy
	err := r.db.Where("site_id = ?", siteID).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新策略
			m := ToModel(policy)
			m.SiteID = &siteID
			if err := r.db.Create(m).Error; err != nil {
				return nil, err
			}
			return ToDomain(m), nil
		}
		return nil, err
	}

	// 更新现有策略
	existing.TimeWindowSeconds = policy.TimeWindowSeconds
	existing.RequestThreshold = policy.RequestThreshold
	existing.DifficultyIncrement = policy.DifficultyIncrement
	existing.BatchIncrement = policy.BatchIncrement
	existing.FailureWeight = policy.FailureWeight
	existing.Enabled = policy.Enabled

	if err := r.db.Save(&existing).Error; err != nil {
		return nil, err
	}
	return ToDomain(&existing), nil
}

// Ensure compile-time interface compliance
var _ ippolicy.Repository = (*IpPolicyRepo)(nil)