package challenge

import "github.com/google/uuid"

// Repository 挑战仓库接口
type Repository interface {
	FindByID(id uuid.UUID) (*Challenge, error)
	FindByIDWithSite(id uuid.UUID) (*Challenge, error)
	FindByBatchID(batchID uuid.UUID) ([]*Challenge, error)
	FindBySiteIDAndToken(siteID uint, token string) (*Challenge, error)
	Create(challenge *Challenge) error
	Save(challenge *Challenge) error

	// 统计
	CountBySiteID(siteID uint) (int64, error)
	CountSolvedBySiteID(siteID uint) (int64, error)
	CountBySiteIDSince(siteID uint, since string) (int64, error)
	CountSolvedBySiteIDSince(siteID uint, since string) (int64, error)
}