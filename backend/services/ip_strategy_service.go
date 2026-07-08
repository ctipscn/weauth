package services

import (
	"github.com/ctipscn/weauth/backend/database"
	"github.com/ctipscn/weauth/backend/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

// GetIpPolicy 获取站点的 IP 策略（包含权限检查）
func GetIpPolicy(siteID uint, userID uint) (*models.IpPolicy, error) {
	// 检查站点是否属于该用户
	var site models.Site
	if err := database.DB.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}

	return GetIpPolicyInternal(siteID)
}

// GetIpPolicyInternal 获取站点的 IP 策略（内部使用，不检查权限）
func GetIpPolicyInternal(siteID uint) (*models.IpPolicy, error) {
	var policy models.IpPolicy

	// 先查找站点特定策略
	err := database.DB.Where("site_id = ?", siteID).First(&policy).Error
	if err == nil {
		return &policy, nil
	}

	// 如果没有站点特定策略，查找全局策略
	err = database.DB.Where("site_id IS NULL").First(&policy).Error
	if err == nil {
		return &policy, nil
	}

	// 没有配置策略，返回默认值
	return &models.IpPolicy{
		TimeWindowSeconds:   300,
		RequestThreshold:    10,
		DifficultyIncrement: 1,
		FailureWeight:       2.0,
		Enabled:             true,
	}, nil
}

// UpdateIpPolicy 更新或创建 IP 策略（包含权限检查）
func UpdateIpPolicy(siteID uint, userID uint, updates map[string]interface{}) (*models.IpPolicy, error) {
	// 检查站点是否属于该用户
	var site models.Site
	if err := database.DB.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}

	var policy models.IpPolicy

	err := database.DB.Where("site_id = ?", siteID).First(&policy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新策略
			policy = models.IpPolicy{
				SiteID: &siteID,
			}
			database.DB.Create(&policy)
		} else {
			return nil, err
		}
	}

	if err := database.DB.Model(&policy).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &policy, nil
}

// GetOrCreateIpRecord 获取或创建 IP 记录
func GetOrCreateIpRecord(siteID uint, clientIP string, policy *models.IpPolicy) (*models.IpRecord, error) {
	var record models.IpRecord

	err := database.DB.Where("site_id = ? AND client_ip = ?", siteID, clientIP).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新记录
			record = models.IpRecord{
				SiteID:      siteID,
				ClientIP:    clientIP,
				WindowStart: time.Now(),
			}
			if err := database.DB.Create(&record).Error; err != nil {
				return nil, err
			}
			return &record, nil
		}
		return nil, err
	}

	// 检查时间窗口是否过期
	if record.IsWindowExpired(policy.TimeWindowSeconds) {
		record.ResetWindow()
		if err := database.DB.Save(&record).Error; err != nil {
			return nil, err
		}
	}

	return &record, nil
}

// CalculateDifficulty 计算动态难度
func CalculateDifficulty(site *models.Site, clientIP string) (int, error) {
	policy, err := GetIpPolicyInternal(site.ID)
	if err != nil {
		return site.BaseDifficulty, nil
	}

	if !policy.Enabled {
		return site.BaseDifficulty, nil
	}

	record, err := GetOrCreateIpRecord(site.ID, clientIP, policy)
	if err != nil {
		return site.BaseDifficulty, nil
	}

	// 计算有效难度
	difficulty := policy.GetEffectiveDifficulty(
		site.BaseDifficulty,
		site.MaxDifficulty,
		record.RequestCount,
		record.FailureCount,
	)

	return difficulty, nil
}

// CalculateBatchCount 根据风控值计算批量模式需要的 PoW 数量
func CalculateBatchCount(site *models.Site, clientIP string) (int, error) {
	policy, err := GetIpPolicyInternal(site.ID)
	if err != nil {
		return site.MinBatchCount, nil
	}

	if !policy.Enabled {
		return site.MinBatchCount, nil
	}

	record, err := GetOrCreateIpRecord(site.ID, clientIP, policy)
	if err != nil {
		return site.MinBatchCount, nil
	}

	// 使用 IpPolicy 的方法计算有效批次数量
	batchCount := policy.GetEffectiveBatchCount(
		site.MinBatchCount,
		site.MaxBatchCount,
		record.RequestCount,
		record.FailureCount,
	)

	return batchCount, nil
}

// IncrementRequestCount 增加请求计数
func IncrementRequestCount(siteID uint, clientIP string) error {
	policy, _ := GetIpPolicyInternal(siteID)
	record, err := GetOrCreateIpRecord(siteID, clientIP, policy)
	if err != nil {
		return err
	}

	record.IncrementRequest()
	return database.DB.Save(record).Error
}

// IncrementFailureCount 增加失败计数
func IncrementFailureCount(siteID uint, clientIP string) error {
	policy, _ := GetIpPolicyInternal(siteID)
	record, err := GetOrCreateIpRecord(siteID, clientIP, policy)
	if err != nil {
		return err
	}

	record.IncrementFailure()
	return database.DB.Save(record).Error
}

// ListIpRecords 获取 IP 记录列表（包含权限检查）
func ListIpRecords(siteID uint, userID uint) ([]models.IpRecord, error) {
	// 检查站点是否属于该用户
	var site models.Site
	if err := database.DB.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}

	var records []models.IpRecord
	err := database.DB.Where("site_id = ?", siteID).
		Order("last_request_at DESC").Find(&records).Error
	return records, err
}

// DeleteIpRecord 删除 IP 记录（包含权限检查）
func DeleteIpRecord(siteID uint, clientIP string, userID uint) error {
	// 检查站点是否属于该用户
	var site models.Site
	if err := database.DB.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSiteNotFound
		}
		return err
	}

	return database.DB.Where("site_id = ? AND client_ip = ?", siteID, clientIP).
		Delete(&models.IpRecord{}).Error
}
