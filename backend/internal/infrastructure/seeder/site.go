package seeder

import (
	"strings"

	"github.com/ctipscn/weauth/backend/internal/domain/site"
	domainuser "github.com/ctipscn/weauth/backend/internal/domain/user"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
)

// SeedSiteOption 为指定用户生成 site + domain 配置的参数。
type SeedSiteOption struct {
	Name             string
	Domains          []string
	BaseDifficulty   int
	MaxDifficulty    int
	ChallengeTimeout int
	TokenTimeout     int
	Enabled          bool
	BatchModeEnabled bool
	BatchDifficulty  int
	MinBatchCount    int
	MaxBatchCount    int
}

// SeedSiteForUser 使用领域实体为指定用户创建站点和域名。
func SeedSiteForUser(
	siteRepo site.Repository,
	user *domainuser.User,
	opt SeedSiteOption,
	log *infralogger.LogService,
) (*site.Site, error) {

	log.Infof(
		"[SEED] SeedSiteForUser: creating site (userID=%d, username=%s, name=%s, domains=%v)",
		user.ID,
		user.Username,
		opt.Name,
		opt.Domains,
	)

	s := site.NewSite(user.ID, opt.Name)
	s.BaseDifficulty = opt.BaseDifficulty
	s.MaxDifficulty = opt.MaxDifficulty
	s.ChallengeTimeout = opt.ChallengeTimeout
	s.TokenTimeout = opt.TokenTimeout
	s.Enabled = opt.Enabled
	s.BatchModeEnabled = opt.BatchModeEnabled
	s.BatchDifficulty = opt.BatchDifficulty
	s.MinBatchCount = opt.MinBatchCount
	s.MaxBatchCount = opt.MaxBatchCount

	err := siteRepo.Transaction(func(repo site.Repository) error {
		if err := repo.Create(s); err != nil {
			return err
		}

		for _, d := range opt.Domains {
			d = strings.TrimSpace(d)
			if d != "" {
				if _, err := repo.AddDomain(s.ID, d); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Errorf("[SEED] SeedSiteForUser: failed to create site: %v", err)
		return nil, err
	}

	// 重新查询以获取包含域名关联数据的完整站点
	result, err := siteRepo.FindByID(s.ID)
	if err != nil {
		log.Errorf("[SEED] SeedSiteForUser: failed to reload site (siteID=%d): %v", s.ID, err)
		return nil, err
	}

	log.Infof(
		"[SEED] SeedSiteForUser: site created (siteID=%d, name=%s, userID=%d, domains=%v)",
		result.ID, opt.Name, user.ID, opt.Domains,
	)
	return result, nil
}
