package verifylog

import (
	"github.com/ctipscn/weauth/backend/internal/domain/verifylog"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
	"gorm.io/gorm"
)

// VerifyLogRepo 验证日志仓储实现
type VerifyLogRepo struct {
	db *gorm.DB
}

// NewVerifyLogRepo 创建验证日志仓储
func NewVerifyLogRepo(driver *database.PostgresDriver) *VerifyLogRepo {
	return &VerifyLogRepo{db: driver.DB()}
}

// Create 创建验证日志
func (r *VerifyLogRepo) Create(log *verifylog.VerificationLog) error {
	m := ToModel(log)
	return r.db.Create(m).Error
}

// CountBySiteID 统计站点的验证总数
func (r *VerifyLogRepo) CountBySiteID(siteID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.VerificationLog{}).Where("site_id = ?", siteID).Count(&count).Error
	return count, err
}

// CountSuccessBySiteID 统计站点成功的验证数
func (r *VerifyLogRepo) CountSuccessBySiteID(siteID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.VerificationLog{}).Where("site_id = ? AND success = true", siteID).Count(&count).Error
	return count, err
}

// Ensure compile-time interface compliance
var _ verifylog.Repository = (*VerifyLogRepo)(nil)