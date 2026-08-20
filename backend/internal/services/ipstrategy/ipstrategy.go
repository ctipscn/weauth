package ipstrategy

import (
	"errors"
	"time"

	"github.com/ctipscn/weauth/backend/internal/database"
	"github.com/ctipscn/weauth/backend/internal/models"
	"github.com/ctipscn/weauth/backend/internal/services/website"

	"gorm.io/gorm"
)

// IPStrategyService IP 策略服务，封装动态难度、批量计数与 IP 记录管理。
type IPStrategyService struct {
	store database.Storer
}

// NewService 创建 IP 策略服务实例。
func NewService(store database.Storer) *IPStrategyService {
	return &IPStrategyService{store: store}
}

// GetIpPolicy 获取站点的 IP 策略（包含权限检查）
func (s *IPStrategyService) GetIpPolicy(siteID uint, userID uint) (*models.IpPolicy, error) {
	db := s.store.DB()
	var site models.Site
	if err := db.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, website.ErrSiteNotFound
		}
		return nil, err
	}

	return s.GetIpPolicyInternal(siteID)
}

// GetIpPolicyInternal 获取站点的 IP 策略（内部使用，不检查权限）
func (s *IPStrategyService) GetIpPolicyInternal(siteID uint) (*models.IpPolicy, error) {
	db := s.store.DB()
	var policy models.IpPolicy

	err := db.Where("site_id = ?", siteID).First(&policy).Error
	if err == nil {
		return &policy, nil
	}

	err = db.Where("site_id IS NULL").First(&policy).Error
	if err == nil {
		return &policy, nil
	}

	return &models.IpPolicy{
		TimeWindowSeconds:   300,
		RequestThreshold:    10,
		DifficultyIncrement: 1,
		FailureWeight:       2.0,
		Enabled:             true,
	}, nil
}

// UpdateIpPolicy 更新或创建 IP 策略（包含权限检查）
func (s *IPStrategyService) UpdateIpPolicy(siteID uint, userID uint, updates map[string]any) (*models.IpPolicy, error) {
	db := s.store.DB()
	var site models.Site
	if err := db.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, website.ErrSiteNotFound
		}
		return nil, err
	}

	var policy models.IpPolicy

	err := db.Where("site_id = ?", siteID).First(&policy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			policy = models.IpPolicy{
				SiteID: &siteID,
			}
			db.Create(&policy)
		} else {
			return nil, err
		}
	}

	if err := db.Model(&policy).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &policy, nil
}

// GetOrCreateIpRecord 获取或创建 IP 记录
func (s *IPStrategyService) GetOrCreateIpRecord(siteID uint, clientIP string, policy *models.IpPolicy) (*models.IpRecord, error) {
	db := s.store.DB()
	var record models.IpRecord

	err := db.Where("site_id = ? AND client_ip = ?", siteID, clientIP).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			record = models.IpRecord{
				SiteID:      siteID,
				ClientIP:    clientIP,
				WindowStart: time.Now(),
			}
			if err := db.Create(&record).Error; err != nil {
				return nil, err
			}
			return &record, nil
		}
		return nil, err
	}

	if record.IsWindowExpired(policy.TimeWindowSeconds) {
		record.ResetWindow()
		if err := db.Save(&record).Error; err != nil {
			return nil, err
		}
	}

	return &record, nil
}

// CalculateDifficulty 计算动态难度
func (s *IPStrategyService) CalculateDifficulty(site *models.Site, clientIP string) (int, error) {
	policy, err := s.GetIpPolicyInternal(site.ID)
	if err != nil {
		return site.BaseDifficulty, nil
	}

	if !policy.Enabled {
		return site.BaseDifficulty, nil
	}

	record, err := s.GetOrCreateIpRecord(site.ID, clientIP, policy)
	if err != nil {
		return site.BaseDifficulty, nil
	}

	difficulty := policy.GetEffectiveDifficulty(
		site.BaseDifficulty,
		site.MaxDifficulty,
		record.RequestCount,
		record.FailureCount,
	)

	return difficulty, nil
}

// CalculateBatchCount 根据风控值计算批量模式需要的 PoW 数量
func (s *IPStrategyService) CalculateBatchCount(site *models.Site, clientIP string) (int, error) {
	policy, err := s.GetIpPolicyInternal(site.ID)
	if err != nil {
		return site.MinBatchCount, nil
	}

	if !policy.Enabled {
		return site.MinBatchCount, nil
	}

	record, err := s.GetOrCreateIpRecord(site.ID, clientIP, policy)
	if err != nil {
		return site.MinBatchCount, nil
	}

	batchCount := policy.GetEffectiveBatchCount(
		site.MinBatchCount,
		site.MaxBatchCount,
		record.RequestCount,
		record.FailureCount,
	)

	return batchCount, nil
}

// IncrementRequestCount 增加请求计数
func (s *IPStrategyService) IncrementRequestCount(siteID uint, clientIP string) error {
	db := s.store.DB()
	policy, _ := s.GetIpPolicyInternal(siteID)
	record, err := s.GetOrCreateIpRecord(siteID, clientIP, policy)
	if err != nil {
		return err
	}

	record.IncrementRequest()
	return db.Save(record).Error
}

// IncrementFailureCount 增加失败计数
func (s *IPStrategyService) IncrementFailureCount(siteID uint, clientIP string) error {
	db := s.store.DB()
	policy, _ := s.GetIpPolicyInternal(siteID)
	record, err := s.GetOrCreateIpRecord(siteID, clientIP, policy)
	if err != nil {
		return err
	}

	record.IncrementFailure()
	return db.Save(record).Error
}

// ListIpRecords 获取 IP 记录列表（包含权限检查）
func (s *IPStrategyService) ListIpRecords(siteID uint, userID uint) ([]models.IpRecord, error) {
	db := s.store.DB()
	var site models.Site
	if err := db.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, website.ErrSiteNotFound
		}
		return nil, err
	}

	var records []models.IpRecord
	err := db.Where("site_id = ?", siteID).
		Order("last_request_at DESC").Find(&records).Error
	return records, err
}

// DeleteIpRecord 删除 IP 记录（包含权限检查）
func (s *IPStrategyService) DeleteIpRecord(siteID uint, clientIP string, userID uint) error {
	db := s.store.DB()
	var site models.Site
	if err := db.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return website.ErrSiteNotFound
		}
		return err
	}

	return db.Where("site_id = ? AND client_ip = ?", siteID, clientIP).
		Delete(&models.IpRecord{}).Error
}
