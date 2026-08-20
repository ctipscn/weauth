package website

import (
	"errors"
	"strings"

	"github.com/ctipscn/weauth/backend/internal/database"
	"github.com/ctipscn/weauth/backend/internal/models"

	"gorm.io/gorm"
)

// SiteService 站点管理服务，封装站点 CRUD、域名验证与 secret 管理。
type SiteService struct {
	store database.Storer
}

// NewSiteService 创建 SiteService 实例。
func NewSiteService(store database.Storer) *SiteService {
	return &SiteService{store: store}
}

// GetSiteBySitekey 根据 sitekey 获取站点
func (s *SiteService) GetSiteBySitekey(sitekey string) (*models.Site, error) {
	db := s.store.DB()
	var site models.Site
	err := db.Preload("Domains").Preload("IpPolicy").
		Where("sitekey = ?", sitekey).First(&site).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}

	if !site.Enabled {
		return nil, ErrSiteDisabled
	}

	return &site, nil
}

// GetSiteBySecret 根据 secret 获取站点
func (s *SiteService) GetSiteBySecret(secret string) (*models.Site, error) {
	db := s.store.DB()
	var site models.Site
	err := db.Where("secret = ?", secret).First(&site).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}

	if !site.Enabled {
		return nil, ErrSiteDisabled
	}

	return &site, nil
}

// ValidateDomain 验证域名是否在白名单中
func (s *SiteService) ValidateDomain(site *models.Site, origin string) bool {
	if len(site.Domains) == 0 {
		return true
	}

	hostname := extractHostname(origin)
	if hostname == "" {
		return false
	}

	for _, domain := range site.Domains {
		if matchDomain(domain.Domain, hostname) {
			return true
		}
	}

	return false
}

// extractHostname 从 origin 提取 hostname
func extractHostname(origin string) string {
	origin = strings.TrimPrefix(origin, "http://")
	origin = strings.TrimPrefix(origin, "https://")

	if idx := strings.Index(origin, ":"); idx != -1 {
		origin = origin[:idx]
	}

	if idx := strings.Index(origin, "/"); idx != -1 {
		origin = origin[:idx]
	}

	return strings.ToLower(origin)
}

// matchDomain 匹配域名（支持通配符 *.example.com）
func matchDomain(pattern, hostname string) bool {
	pattern = strings.ToLower(pattern)
	hostname = strings.ToLower(hostname)

	if pattern == hostname {
		return true
	}

	if strings.HasPrefix(pattern, "*.") {
		suffix := pattern[1:]
		return strings.HasSuffix(hostname, suffix)
	}

	return false
}

// CreateSite 创建站点
func (s *SiteService) CreateSite(userID uint, name string, domains []string, opts ...CreateSiteOptions) (*models.Site, error) {
	db := s.store.DB()
	site := &models.Site{
		UserID: userID,
		Name:   name,
		// 默认启用批量模式
		BatchModeEnabled: true,
		BatchDifficulty:  4,
		MinBatchCount:    1,
		MaxBatchCount:    10,
	}

	if len(opts) > 0 {
		opt := opts[0]
		if opt.BatchModeEnabled != nil {
			site.BatchModeEnabled = *opt.BatchModeEnabled
		}
		if opt.BatchDifficulty != nil {
			site.BatchDifficulty = *opt.BatchDifficulty
		}
		if opt.MinBatchCount != nil {
			site.MinBatchCount = *opt.MinBatchCount
		}
		if opt.MaxBatchCount != nil {
			site.MaxBatchCount = *opt.MaxBatchCount
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(site).Error; err != nil {
			return err
		}

		for _, domain := range domains {
			if domain = strings.TrimSpace(domain); domain != "" {
				siteDomain := &models.SiteDomain{
					SiteID: site.ID,
					Domain: domain,
				}
				if err := tx.Create(siteDomain).Error; err != nil {
					return err
				}
			}
		}

		ipPolicy := &models.IpPolicy{
			SiteID: &site.ID,
		}
		if err := tx.Create(ipPolicy).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetSiteByIDInternal(site.ID)
}

// GetSiteByID 根据 ID 获取站点（包含权限检查）
func (s *SiteService) GetSiteByID(id uint, userID uint) (*models.Site, error) {
	db := s.store.DB()
	var site models.Site
	err := db.Preload("Domains").Preload("IpPolicy").
		Where("id = ? AND user_id = ?", id, userID).
		First(&site).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}
	return &site, nil
}

// GetSiteByIDInternal 内部使用，不检查权限
func (s *SiteService) GetSiteByIDInternal(id uint) (*models.Site, error) {
	db := s.store.DB()
	var site models.Site
	err := db.Preload("Domains").Preload("IpPolicy").
		First(&site, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}
	return &site, nil
}

// ListSites 获取指定用户的站点列表
func (s *SiteService) ListSites(userID uint) ([]models.Site, error) {
	db := s.store.DB()
	var sites []models.Site
	err := db.Preload("Domains").Preload("IpPolicy").
		Where("user_id = ?", userID).
		Order("created_at DESC").Find(&sites).Error
	return sites, err
}

// UpdateSite 更新站点（包含权限检查）
func (s *SiteService) UpdateSite(id uint, userID uint, updates map[string]any) (*models.Site, error) {
	db := s.store.DB()
	var site models.Site
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}

	if err := db.Model(&site).Updates(updates).Error; err != nil {
		return nil, err
	}

	return s.GetSiteByID(id, userID)
}

// DeleteSite 删除站点（包含权限检查）
func (s *SiteService) DeleteSite(id uint, userID uint) error {
	db := s.store.DB()
	result := db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Site{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSiteNotFound
	}
	return nil
}

// AddDomain 添加域名（包含权限检查）
func (s *SiteService) AddDomain(siteID uint, userID uint, domain string) (*models.SiteDomain, error) {
	db := s.store.DB()
	var site models.Site
	if err := db.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}

	siteDomain := &models.SiteDomain{
		SiteID: siteID,
		Domain: strings.TrimSpace(domain),
	}

	if err := db.Create(siteDomain).Error; err != nil {
		return nil, err
	}

	return siteDomain, nil
}

// DeleteDomain 删除域名（包含权限检查）
func (s *SiteService) DeleteDomain(siteID, domainID uint, userID uint) error {
	db := s.store.DB()
	var site models.Site
	if err := db.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSiteNotFound
		}
		return err
	}

	return db.Where("id = ? AND site_id = ?", domainID, siteID).
		Delete(&models.SiteDomain{}).Error
}

// RegenerateSiteSecret 重新生成站点 secret（包含权限检查）
func (s *SiteService) RegenerateSiteSecret(id uint, userID uint) (string, error) {
	db := s.store.DB()
	var site models.Site
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrSiteNotFound
		}
		return "", err
	}

	site.RegenerateSecret()

	if err := db.Save(&site).Error; err != nil {
		return "", err
	}

	return site.Secret, nil
}
