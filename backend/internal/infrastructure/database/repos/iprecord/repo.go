package iprecord

import (
	"github.com/ctipscn/weauth/backend/internal/domain"
	"github.com/ctipscn/weauth/backend/internal/domain/iprecord"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
	"gorm.io/gorm"
)

// IpRecordRepo IP 记录仓储实现
type IpRecordRepo struct {
	db *gorm.DB
}

// NewIpRecordRepo 创建 IP 记录仓储
func NewIpRecordRepo(driver *database.PostgresDriver) *IpRecordRepo {
	return &IpRecordRepo{db: driver.DB()}
}

// FindBySiteIDAndIP 根据站点ID和IP查找记录
func (r *IpRecordRepo) FindBySiteIDAndIP(siteID uint, ip string) (*iprecord.IpRecord, error) {
	var m models.IpRecord
	if err := r.db.Where("site_id = ? AND client_ip = ?", siteID, ip).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// FindAllBySiteID 根据站点ID查找所有IP记录
func (r *IpRecordRepo) FindAllBySiteID(siteID uint) ([]*iprecord.IpRecord, error) {
	var ms []models.IpRecord
	if err := r.db.Where("site_id = ?", siteID).Order("last_request_at DESC").Find(&ms).Error; err != nil {
		return nil, err
	}
	result := make([]*iprecord.IpRecord, 0, len(ms))
	for i := range ms {
		result = append(result, ToDomain(&ms[i]))
	}
	return result, nil
}

// FindRecentBySiteID 统计站点在时间窗口内的IP数
func (r *IpRecordRepo) FindRecentBySiteID(siteID uint, window int) ([]*iprecord.IpRecord, error) {
	var ms []models.IpRecord
	if err := r.db.Where("site_id = ?", siteID).Find(&ms).Error; err != nil {
		return nil, err
	}
	result := make([]*iprecord.IpRecord, 0, len(ms))
	for i := range ms {
		result = append(result, ToDomain(&ms[i]))
	}
	return result, nil
}

// Create 创建IP记录
func (r *IpRecordRepo) Create(record *iprecord.IpRecord) error {
	m := ToModel(record)
	return r.db.Create(m).Error
}

// Update 更新IP记录
func (r *IpRecordRepo) Update(record *iprecord.IpRecord) error {
	m := ToModel(record)
	return r.db.Save(m).Error
}

// Save 创建或更新IP记录
func (r *IpRecordRepo) Save(record *iprecord.IpRecord) error {
	m := ToModel(record)
	return r.db.Save(m).Error
}

// DeleteBySiteID 删除站点的所有IP记录
func (r *IpRecordRepo) DeleteBySiteID(siteID uint) error {
	return r.db.Where("site_id = ?", siteID).Delete(&models.IpRecord{}).Error
}

// ListBySiteID 根据站点ID列出IP记录
func (r *IpRecordRepo) ListBySiteID(siteID uint) ([]*iprecord.IpRecord, error) {
	return r.FindAllBySiteID(siteID)
}

// DeleteBySiteIDAndIP 根据站点ID和IP删除记录
func (r *IpRecordRepo) DeleteBySiteIDAndIP(siteID uint, clientIP string) error {
	return r.db.Where("site_id = ? AND client_ip = ?", siteID, clientIP).Delete(&models.IpRecord{}).Error
}

// CountBySiteID 统计站点的IP记录总数
func (r *IpRecordRepo) CountBySiteID(siteID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.IpRecord{}).Where("site_id = ?", siteID).Count(&count).Error
	return count, err
}

// CountRecentBySiteID 统计站点在时间窗口内的活跃IP数
func (r *IpRecordRepo) CountRecentBySiteID(siteID uint, window int) (int64, error) {
	var count int64
	err := r.db.Model(&models.IpRecord{}).
		Where("site_id = ?", siteID).
		Count(&count).Error
	return count, err
}

// Ensure compile-time interface compliance
var _ iprecord.Repository = (*IpRecordRepo)(nil)