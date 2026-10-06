package challenge

import (
	"errors"

	"github.com/ctipscn/weauth/backend/internal/domain"
	"github.com/ctipscn/weauth/backend/internal/domain/challenge"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChallengeRepo 挑战仓储实现
type ChallengeRepo struct {
	db *gorm.DB
}

// NewChallengeRepo 创建挑战仓储
func NewChallengeRepo(driver *database.PostgresDriver) *ChallengeRepo {
	return &ChallengeRepo{db: driver.DB()}
}

// FindByID 根据ID查找挑战
func (r *ChallengeRepo) FindByID(id uuid.UUID) (*challenge.Challenge, error) {
	var m models.Challenge
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrChallengeNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// FindByIDWithSite 根据ID查找挑战（预加载站点）
func (r *ChallengeRepo) FindByIDWithSite(id uuid.UUID) (*challenge.Challenge, error) {
	var m models.Challenge
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrChallengeNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// FindByBatchID 根据批次ID查找所有挑战
func (r *ChallengeRepo) FindByBatchID(batchID uuid.UUID) ([]*challenge.Challenge, error) {
	var ms []models.Challenge
	if err := r.db.Where("batch_id = ?", batchID).Order("batch_index").Find(&ms).Error; err != nil {
		return nil, err
	}
	result := make([]*challenge.Challenge, 0, len(ms))
	for i := range ms {
		result = append(result, ToDomain(&ms[i]))
	}
	return result, nil
}

// FindBySiteIDAndToken 根据站点ID和token查找挑战
func (r *ChallengeRepo) FindBySiteIDAndToken(siteID uint, token string) (*challenge.Challenge, error) {
	var m models.Challenge
	if err := r.db.Where("site_id = ? AND token = ?", siteID, token).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrChallengeNotFound
		}
		return nil, err
	}
	return ToDomain(&m), nil
}

// Create 创建挑战
func (r *ChallengeRepo) Create(ch *challenge.Challenge) error {
	m := ToModel(ch)
	return r.db.Create(m).Error
}

// Save 保存挑战
func (r *ChallengeRepo) Save(ch *challenge.Challenge) error {
	m := ToModel(ch)
	return r.db.Save(m).Error
}

// CountBySiteID 统计站点的挑战总数
func (r *ChallengeRepo) CountBySiteID(siteID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Challenge{}).Where("site_id = ?", siteID).Count(&count).Error
	return count, err
}

// CountSolvedBySiteID 统计站点已解决的挑战数
func (r *ChallengeRepo) CountSolvedBySiteID(siteID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Challenge{}).Where("site_id = ? AND solved = true", siteID).Count(&count).Error
	return count, err
}

// CountBySiteIDSince 统计站点在指定时间后的挑战数
func (r *ChallengeRepo) CountBySiteIDSince(siteID uint, since string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Challenge{}).
		Where("site_id = ? AND created_at > ?", siteID, since).
		Count(&count).Error
	return count, err
}

// CountSolvedBySiteIDSince 统计站点在指定时间后已解决的挑战数
func (r *ChallengeRepo) CountSolvedBySiteIDSince(siteID uint, since string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Challenge{}).
		Where("site_id = ? AND solved = true AND created_at > ?", siteID, since).
		Count(&count).Error
	return count, err
}

// Ensure compile-time interface compliance
var _ challenge.Repository = (*ChallengeRepo)(nil)