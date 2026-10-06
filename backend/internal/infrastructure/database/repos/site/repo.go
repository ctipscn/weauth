package site

import (
	"errors"

	"github.com/ctipscn/weauth/backend/internal/domain"
	"github.com/ctipscn/weauth/backend/internal/domain/site"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
	"gorm.io/gorm"
)

// SiteRepo 站点仓储实现
type SiteRepo struct {
	db *gorm.DB
}

// NewSiteRepo 创建站点仓储
func NewSiteRepo(driver *database.PostgresDriver) *SiteRepo {
	return &SiteRepo{db: driver.DB()}
}

// preload 返回带 Domains 预加载的 DB 实例
func (r *SiteRepo) preload() *gorm.DB {
	return r.db.Preload("Domains")
}

// FindByID 根据ID查找站点
func (r *SiteRepo) FindByID(id uint) (*site.Site, error) {
	var m models.Site
	if err := r.preload().First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSiteNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// FindByIDWithUser 根据ID查找站点（包含权限校验）
func (r *SiteRepo) FindByIDWithUser(id uint, userID uint) (*site.Site, error) {
	var m models.Site
	if err := r.preload().Where("id = ? AND user_id = ?", id, userID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSiteNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// FindBySitekey 根据sitekey查找站点
func (r *SiteRepo) FindBySitekey(sitekey string) (*site.Site, error) {
	var m models.Site
	if err := r.preload().Where("sitekey = ?", sitekey).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSiteNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// FindBySecret 根据secret查找站点
func (r *SiteRepo) FindBySecret(secret string) (*site.Site, error) {
	var m models.Site
	if err := r.preload().Where("secret = ?", secret).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSiteNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// FindAll 获取所有站点列表
func (r *SiteRepo) FindAll() ([]*site.Site, error) {
	var ms []models.Site
	if err := r.preload().Order("created_at DESC").Find(&ms).Error; err != nil {
		return nil, err
	}
	result := make([]*site.Site, 0, len(ms))
	for i := range ms {
		result = append(result, ToDomain(&ms[i]))
	}
	return result, nil
}

// ListByUser 获取指定用户的站点列表
func (r *SiteRepo) ListByUser(userID uint) ([]*site.Site, error) {
	var ms []models.Site
	if err := r.preload().Where("user_id = ?", userID).Order("created_at DESC").Find(&ms).Error; err != nil {
		return nil, err
	}
	result := make([]*site.Site, 0, len(ms))
	for i := range ms {
		result = append(result, ToDomain(&ms[i]))
	}
	return result, nil
}

// Create 创建站点
func (r *SiteRepo) Create(s *site.Site) error {
	m := ToModel(s)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	s.ID = m.ID
	s.Sitekey = m.Sitekey
	s.Secret = m.Secret
	return nil
}

// Save 全量更新站点
func (r *SiteRepo) Save(s *site.Site) error {
	m := ToModel(s)
	return r.db.Save(m).Error
}

// AddDomain 添加域名（拒绝重复添加）
func (r *SiteRepo) AddDomain(siteID uint, domainName string) (*site.SiteDomain, error) {
	// 检查是否已存在
	var count int64
	if err := r.db.Model(&models.Domain{}).Where("site_id = ? AND domain = ?", siteID, domainName).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, domain.ErrDomainAlreadyExists
	}

	m := &models.Domain{
		SiteID: siteID,
		Domain: domainName,
	}
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return &site.SiteDomain{
		ID:        m.ID,
		SiteID:    m.SiteID,
		Domain:    m.Domain,
		CreatedAt: m.CreatedAt,
	}, nil
}

// DeleteDomain 删除域名
func (r *SiteRepo) DeleteDomain(siteID, domainID uint) error {
	result := r.db.Where("id = ? AND site_id = ?", domainID, siteID).Delete(&models.Domain{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrDomainNotFound
	}
	return nil
}

// Transaction 在事务内执行操作
func (r *SiteRepo) Transaction(fn func(repo site.Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &SiteRepo{db: tx}
		return fn(txRepo)
	})
}

// Delete 删除站点（同时清理关联域名，防止软删除遗留脏数据）
func (r *SiteRepo) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先清理关联域名（domains 表无 deleted_at，需硬删）
		if err := tx.Where("site_id = ?", id).Delete(&models.Domain{}).Error; err != nil {
			return err
		}

		// 软删除站点本身
		result := tx.Delete(&models.Site{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domain.ErrSiteNotFound
		}
		return nil
	})
}

// FindAllUnscoped 获取所有站点列表（包含已软删除的记录）
func (r *SiteRepo) FindAllUnscoped() ([]*site.Site, error) {
	var ms []models.Site
	if err := r.db.Unscoped().Preload("Domains").Order("created_at DESC").Find(&ms).Error; err != nil {
		return nil, err
	}
	result := make([]*site.Site, 0, len(ms))
	for i := range ms {
		result = append(result, ToDomain(&ms[i]))
	}
	return result, nil
}

// FindByIDUnscoped 根据 ID 查找站点（包含已软删除的记录）
func (r *SiteRepo) FindByIDUnscoped(id uint) (*site.Site, error) {
	var m models.Site
	if err := r.db.Unscoped().Preload("Domains").First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSiteNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// Ensure compile-time interface compliance
var _ site.Repository = (*SiteRepo)(nil)