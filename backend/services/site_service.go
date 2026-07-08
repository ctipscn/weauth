package services

import (
	"github.com/ctipscn/weauth/backend/database"
	"github.com/ctipscn/weauth/backend/models"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrSiteNotFound    = errors.New("site not found")
	ErrSiteDisabled    = errors.New("site is disabled")
	ErrInvalidSitekey  = errors.New("invalid sitekey")
	ErrInvalidSecret   = errors.New("invalid secret")
	ErrDomainNotFound  = errors.New("domain not found")
	ErrDomainMismatch  = errors.New("domain not in whitelist")
	ErrPermissionDenied = errors.New("permission denied")
)

// GetSiteBySitekey 根据 sitekey 获取站点
func GetSiteBySitekey(sitekey string) (*models.Site, error) {
	var site models.Site
	err := database.DB.Preload("Domains").Preload("IpPolicy").
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
func GetSiteBySecret(secret string) (*models.Site, error) {
	var site models.Site
	err := database.DB.Where("secret = ?", secret).First(&site).Error
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
func ValidateDomain(site *models.Site, origin string) bool {
	if len(site.Domains) == 0 {
		// 没有配置域名白名单，允许所有
		return true
	}

	// 从 origin 提取 hostname
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
	// 移除协议前缀
	origin = strings.TrimPrefix(origin, "http://")
	origin = strings.TrimPrefix(origin, "https://")

	// 移除端口
	if idx := strings.Index(origin, ":"); idx != -1 {
		origin = origin[:idx]
	}

	// 移除路径
	if idx := strings.Index(origin, "/"); idx != -1 {
		origin = origin[:idx]
	}

	return strings.ToLower(origin)
}

// matchDomain 匹配域名（支持通配符 *.example.com）
func matchDomain(pattern, hostname string) bool {
	pattern = strings.ToLower(pattern)
	hostname = strings.ToLower(hostname)

	// 精确匹配
	if pattern == hostname {
		return true
	}

	// 通配符匹配 *.example.com
	if strings.HasPrefix(pattern, "*.") {
		suffix := pattern[1:] // .example.com
		return strings.HasSuffix(hostname, suffix)
	}

	return false
}

// CreateSiteOptions 创建站点的可选配置
type CreateSiteOptions struct {
	BatchModeEnabled *bool
	BatchDifficulty  *int
	MinBatchCount    *int
	MaxBatchCount    *int
}

// CreateSite 创建站点
func CreateSite(userID uint, name string, domains []string, opts ...CreateSiteOptions) (*models.Site, error) {
	site := &models.Site{
		UserID: userID,
		Name:   name,
		// 默认启用批量模式
		BatchModeEnabled: true,
		BatchDifficulty:  4,
		MinBatchCount:    1,
		MaxBatchCount:    10,
	}

	// 应用可选配置
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

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(site).Error; err != nil {
			return err
		}

		// 添加域名
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

		// 创建默认 IP 策略
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

	// 重新加载以获取关联数据（内部调用，不需要权限检查）
	return GetSiteByIDInternal(site.ID)
}

// GetSiteByID 根据 ID 获取站点（包含权限检查）
func GetSiteByID(id uint, userID uint) (*models.Site, error) {
	var site models.Site
	err := database.DB.Preload("Domains").Preload("IpPolicy").
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
func GetSiteByIDInternal(id uint) (*models.Site, error) {
	var site models.Site
	err := database.DB.Preload("Domains").Preload("IpPolicy").
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
func ListSites(userID uint) ([]models.Site, error) {
	var sites []models.Site
	err := database.DB.Preload("Domains").Preload("IpPolicy").
		Where("user_id = ?", userID).
		Order("created_at DESC").Find(&sites).Error
	return sites, err
}

// UpdateSite 更新站点（包含权限检查）
func UpdateSite(id uint, userID uint, updates map[string]interface{}) (*models.Site, error) {
	var site models.Site
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}

	if err := database.DB.Model(&site).Updates(updates).Error; err != nil {
		return nil, err
	}

	return GetSiteByID(id, userID)
}

// DeleteSite 删除站点（包含权限检查）
func DeleteSite(id uint, userID uint) error {
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Site{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSiteNotFound
	}
	return nil
}

// AddDomain 添加域名（包含权限检查）
func AddDomain(siteID uint, userID uint, domain string) (*models.SiteDomain, error) {
	// 先检查站点是否属于该用户
	var site models.Site
	if err := database.DB.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}

	siteDomain := &models.SiteDomain{
		SiteID: siteID,
		Domain: strings.TrimSpace(domain),
	}

	if err := database.DB.Create(siteDomain).Error; err != nil {
		return nil, err
	}

	return siteDomain, nil
}

// DeleteDomain 删除域名（包含权限检查）
func DeleteDomain(siteID, domainID uint, userID uint) error {
	// 先检查站点是否属于该用户
	var site models.Site
	if err := database.DB.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSiteNotFound
		}
		return err
	}

	return database.DB.Where("id = ? AND site_id = ?", domainID, siteID).
		Delete(&models.SiteDomain{}).Error
}

// RegenerateSiteSecret 重新生成站点 secret（包含权限检查）
func RegenerateSiteSecret(id uint, userID uint) (string, error) {
	var site models.Site
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrSiteNotFound
		}
		return "", err
	}

	site.RegenerateSecret()

	if err := database.DB.Save(&site).Error; err != nil {
		return "", err
	}

	return site.Secret, nil
}
