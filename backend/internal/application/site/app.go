package site

import (
	"errors"
	"strings"

	"github.com/ctipscn/weauth/backend/internal/application/pow"
	"github.com/ctipscn/weauth/backend/internal/domain"
	"github.com/ctipscn/weauth/backend/internal/domain/ippolicy"
	"github.com/ctipscn/weauth/backend/internal/domain/iprecord"
	"github.com/ctipscn/weauth/backend/internal/domain/site"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
)

// SiteApp 站点管理应用服务
type SiteApp struct {
	siteRepo     site.Repository
	ipPolicyRepo ippolicy.Repository
	ipRecordRepo iprecord.Repository
	powApp       *pow.PoWApp
	log          *infralogger.LogService
}

// NewSiteApp 创建站点管理应用服务
func NewSiteApp(
	siteRepo site.Repository,
	ipPolicyRepo ippolicy.Repository,
	ipRecordRepo iprecord.Repository,
	powApp *pow.PoWApp,
	log *infralogger.LogService,
) *SiteApp {
	return &SiteApp{
		siteRepo:     siteRepo,
		ipPolicyRepo: ipPolicyRepo,
		ipRecordRepo: ipRecordRepo,
		powApp:       powApp,
		log:          log,
	}
}

// CreateSite 创建站点（包含域名和默认 IP 策略）
func (app *SiteApp) CreateSite(
	userID uint,
	name string,
	domains []string,
	opts ...CreateSiteOptions,
) (*site.Site, error) {

	s := site.NewSite(userID, name)

	// 应用可选配置
	if len(opts) > 0 {
		opt := opts[0]
		if opt.BatchModeEnabled != nil {
			s.BatchModeEnabled = *opt.BatchModeEnabled
		}
		if opt.BatchDifficulty != nil {
			s.BatchDifficulty = *opt.BatchDifficulty
		}
		if opt.MinBatchCount != nil {
			s.MinBatchCount = *opt.MinBatchCount
		}
		if opt.MaxBatchCount != nil {
			s.MaxBatchCount = *opt.MaxBatchCount
		}
	}

	// 使用事务创建站点、域名和默认 IP 策略
	err := app.siteRepo.Transaction(func(repo site.Repository) error {
		if err := repo.Create(s); err != nil {
			return err
		}

		// 添加域名
		for _, d := range domains {
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
		app.log.Errorf("[SITE] CreateSite failed (userID=%d, name=%s): %v", userID, name, err)
		return nil, err
	}

	result, err := app.siteRepo.FindByID(s.ID)
	if err != nil {
		app.log.Errorf("[SITE] CreateSite reload failed (siteID=%d): %v", s.ID, err)
		return nil, err
	}

	app.log.Infof("[SITE] Site created (siteID=%d, name=%s, userID=%d, domains=%v)", result.ID, name, userID, domains)
	return result, nil
}

// GetSiteByID 根据 ID 获取站点（包含权限检查）
func (app *SiteApp) GetSiteByID(id, userID uint) (*site.Site, error) {
	s, err := app.siteRepo.FindByIDWithUser(id, userID)
	if err != nil {
		app.log.Warnf("[SITE] GetSiteByID failed: site not found (siteID=%d, userID=%d)", id, userID)
		return nil, pow.ErrSiteNotFound
	}
	return s, nil
}

// ListSites 获取指定用户的站点列表
func (app *SiteApp) ListSites(userID uint) ([]*site.Site, error) {
	sites, err := app.siteRepo.ListByUser(userID)
	if err != nil {
		app.log.Errorf("[SITE] ListSites error (userID=%d): %v", userID, err)
		return nil, err
	}
	app.log.Infof("[SITE] Listed sites (userID=%d, count=%d)", userID, len(sites))
	return sites, nil
}

// DeleteSite 删除站点（包含权限检查）
func (app *SiteApp) DeleteSite(id, userID uint) error {
	_, err := app.siteRepo.FindByIDWithUser(id, userID)
	if err != nil {
		app.log.Warnf("[SITE] DeleteSite failed: site not found (siteID=%d, userID=%d)", id, userID)
		return pow.ErrSiteNotFound
	}

	if err := app.siteRepo.Delete(id); err != nil {
		app.log.Errorf("[SITE] DeleteSite error (siteID=%d, userID=%d): %v", id, userID, err)
		return err
	}

	app.log.Infof("[SITE] Site deleted (siteID=%d, userID=%d)", id, userID)
	return nil
}

// UpdateSite 更新站点（包含权限检查）
func (app *SiteApp) UpdateSite(id, userID uint, updates map[string]any) (*site.Site, error) {
	// 验证站点归属
	_, err := app.siteRepo.FindByIDWithUser(id, userID)
	if err != nil {
		app.log.Warnf("[SITE] UpdateSite failed: site not found (siteID=%d, userID=%d)", id, userID)
		return nil, pow.ErrSiteNotFound
	}

	s, err := app.siteRepo.FindByID(id)
	if err != nil {
		app.log.Warnf("[SITE] UpdateSite failed: site not found by ID (siteID=%d)", id)
		return nil, pow.ErrSiteNotFound
	}

	// 应用更新
	if v, ok := updates["name"]; ok {
		s.Name = v.(string)
	}
	if v, ok := updates["base_difficulty"]; ok {
		s.BaseDifficulty = v.(int)
	}
	if v, ok := updates["max_difficulty"]; ok {
		s.MaxDifficulty = v.(int)
	}
	if v, ok := updates["challenge_timeout"]; ok {
		s.ChallengeTimeout = v.(int)
	}
	if v, ok := updates["token_timeout"]; ok {
		s.TokenTimeout = v.(int)
	}
	if v, ok := updates["enabled"]; ok {
		s.Enabled = v.(bool)
	}
	if v, ok := updates["batch_mode_enabled"]; ok {
		s.BatchModeEnabled = v.(bool)
	}
	if v, ok := updates["batch_difficulty"]; ok {
		s.BatchDifficulty = v.(int)
	}
	if v, ok := updates["min_batch_count"]; ok {
		s.MinBatchCount = v.(int)
	}
	if v, ok := updates["max_batch_count"]; ok {
		s.MaxBatchCount = v.(int)
	}

	if err := app.siteRepo.Save(s); err != nil {
		app.log.Errorf("[SITE] UpdateSite save error (siteID=%d, userID=%d): %v", id, userID, err)
		return nil, err
	}

	app.log.Infof("[SITE] Site updated (siteID=%d, userID=%d, updates=%v)", id, userID, updates)
	return s, nil
}

// AddDomain 添加域名（包含权限检查）
func (app *SiteApp) AddDomain(siteID, userID uint, domainName string) (*site.SiteDomain, error) {
	domainName = strings.TrimSpace(domainName)
	_, err := app.siteRepo.FindByIDWithUser(siteID, userID)
	if err != nil {
		app.log.Warnf("[SITE] AddDomain failed: site not found (siteID=%d, userID=%d)", siteID, userID)
		return nil, pow.ErrSiteNotFound
	}

	d, err := app.siteRepo.AddDomain(siteID, domainName)
	if err != nil {
		app.log.Errorf("[SITE] AddDomain error (siteID=%d, domain=%s): %v", siteID, domainName, err)
		return nil, err
	}

	app.log.Infof("[SITE] Domain added (siteID=%d, domainID=%d, domain=%s)", siteID, d.ID, domainName)
	return d, nil
}

// DeleteDomain 删除域名（包含权限检查）
func (app *SiteApp) DeleteDomain(siteID, domainID, userID uint) error {
	_, err := app.siteRepo.FindByIDWithUser(siteID, userID)
	if err != nil {
		app.log.Warnf("[SITE] DeleteDomain failed: site not found (siteID=%d, domainID=%d, userID=%d)", siteID, domainID, userID)
		return pow.ErrSiteNotFound
	}

	err = app.siteRepo.DeleteDomain(siteID, domainID)
	if err != nil {
		app.log.Errorf("[SITE] DeleteDomain error (siteID=%d, domainID=%d): %v", siteID, domainID, err)
		return err
	}

	app.log.Infof("[SITE] Domain deleted (siteID=%d, domainID=%d)", siteID, domainID)
	return nil
}

// RegenerateSiteSecret 重新生成站点 secret（包含权限检查）
func (app *SiteApp) RegenerateSiteSecret(id, userID uint) (string, error) {
	s, err := app.siteRepo.FindByIDWithUser(id, userID)
	if err != nil {
		app.log.Warnf("[SITE] RegenerateSiteSecret failed: site not found (siteID=%d, userID=%d)", id, userID)
		return "", pow.ErrSiteNotFound
	}

	oldSecret := s.Secret[:8] + "..."
	s.RegenerateSecret()
	if err := app.siteRepo.Save(s); err != nil {
		app.log.Errorf("[SITE] RegenerateSiteSecret save error (siteID=%d, userID=%d): %v", id, userID, err)
		return "", err
	}

	app.log.Infof("[SITE] Secret regenerated (siteID=%d, userID=%d, oldPrefix=%s)", id, userID, oldSecret)
	return s.Secret, nil
}

// GetIpPolicy 获取站点的 IP 策略（包含权限检查）
func (app *SiteApp) GetIpPolicy(siteID, userID uint) (*ippolicy.IpPolicy, error) {
	_, err := app.siteRepo.FindByIDWithUser(siteID, userID)
	if err != nil {
		app.log.Warnf("[SITE] GetIpPolicy failed: site not found (siteID=%d, userID=%d)", siteID, userID)
		return nil, pow.ErrSiteNotFound
	}

	policy, err := app.ipPolicyRepo.FindBySiteID(siteID)
	if err != nil {
		// 回退到全局策略
		if errors.Is(err, domain.ErrNotFound) {
			global, globalErr := app.ipPolicyRepo.FindGlobal()
			if globalErr != nil {
				app.log.Warnf("[SITE] GetIpPolicy: falling back to default policy (siteID=%d)", siteID)
				return ippolicy.NewDefaultIpPolicy(), nil
			}
			app.log.Infof("[SITE] GetIpPolicy: using global policy (siteID=%d)", siteID)
			return global, nil
		}
		app.log.Errorf("[SITE] GetIpPolicy error (siteID=%d): %v", siteID, err)
		return nil, err
	}
	return policy, nil
}

// UpdateIpPolicy 更新 IP 策略（包含权限检查）
func (app *SiteApp) UpdateIpPolicy(
	siteID, userID uint,
	updates map[string]any,
) (*ippolicy.IpPolicy, error) {

	_, err := app.siteRepo.FindByIDWithUser(siteID, userID)
	if err != nil {
		app.log.Warnf("[SITE] UpdateIpPolicy failed: site not found (siteID=%d, userID=%d)", siteID, userID)
		return nil, pow.ErrSiteNotFound
	}

	// 获取现有策略或创建新策略
	policy, err := app.ipPolicyRepo.FindBySiteID(siteID)
	if err != nil {
		app.log.Infof("[SITE] UpdateIpPolicy: creating new policy for site (siteID=%d)", siteID)
		policy = ippolicy.NewSiteIpPolicy(siteID)
	}

	// 应用更新
	if v, ok := updates["time_window_seconds"]; ok {
		policy.TimeWindowSeconds = v.(int)
	}
	if v, ok := updates["request_threshold"]; ok {
		policy.RequestThreshold = v.(int)
	}
	if v, ok := updates["difficulty_increment"]; ok {
		policy.DifficultyIncrement = v.(int)
	}
	if v, ok := updates["batch_increment"]; ok {
		policy.BatchIncrement = v.(int)
	}
	if v, ok := updates["failure_weight"]; ok {
		policy.FailureWeight = v.(float64)
	}
	if v, ok := updates["enabled"]; ok {
		policy.Enabled = v.(bool)
	}

	result, err := app.ipPolicyRepo.CreateOrUpdate(siteID, policy)
	if err != nil {
		app.log.Errorf("[SITE] UpdateIpPolicy error (siteID=%d, userID=%d): %v", siteID, userID, err)
		return nil, err
	}

	app.log.Infof("[SITE] IP policy updated (siteID=%d, userID=%d)", siteID, userID)
	return result, nil
}

// ListIpRecords 获取 IP 记录列表（包含权限检查）
func (app *SiteApp) ListIpRecords(siteID, userID uint) ([]*iprecord.IpRecord, error) {
	_, err := app.siteRepo.FindByIDWithUser(siteID, userID)
	if err != nil {
		app.log.Warnf("[SITE] ListIpRecords failed: site not found (siteID=%d, userID=%d)", siteID, userID)
		return nil, pow.ErrSiteNotFound
	}

	records, err := app.ipRecordRepo.ListBySiteID(siteID)
	if err != nil {
		app.log.Errorf("[SITE] ListIpRecords error (siteID=%d): %v", siteID, err)
		return nil, err
	}

	app.log.Infof("[SITE] Listed IP records (siteID=%d, count=%d)", siteID, len(records))
	return records, nil
}

// DeleteIpRecord 删除 IP 记录（包含权限检查）
func (app *SiteApp) DeleteIpRecord(siteID, userID uint, clientIP string) error {
	_, err := app.siteRepo.FindByIDWithUser(siteID, userID)
	if err != nil {
		app.log.Warnf("[SITE] DeleteIpRecord failed: site not found (siteID=%d, userID=%d)", siteID, userID)
		return pow.ErrSiteNotFound
	}

	err = app.ipRecordRepo.DeleteBySiteIDAndIP(siteID, clientIP)
	if err != nil {
		app.log.Infof("[SITE] DeleteIpRecord: record not found or already deleted (siteID=%d, clientIP=%s)", siteID, clientIP)
		return nil // 不存在也算成功
	}

	app.log.Infof("[SITE] IP record deleted (siteID=%d, clientIP=%s)", siteID, clientIP)
	return nil
}
