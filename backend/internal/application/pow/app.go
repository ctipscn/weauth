package pow

import (
	"errors"
	"time"

	"github.com/ctipscn/weauth/backend/internal/domain"
	"github.com/ctipscn/weauth/backend/internal/domain/challenge"
	"github.com/ctipscn/weauth/backend/internal/domain/ippolicy"
	"github.com/ctipscn/weauth/backend/internal/domain/iprecord"
	"github.com/ctipscn/weauth/backend/internal/domain/site"
	"github.com/ctipscn/weauth/backend/internal/domain/verifylog"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"github.com/google/uuid"
)

// PoWApp PoW 应用服务
type PoWApp struct {
	challengeRepo challenge.Repository
	siteRepo      site.Repository
	ipPolicyRepo  ippolicy.Repository
	ipRecordRepo  iprecord.Repository
	verifylogRepo verifylog.Repository
	log           *infralogger.LogService
}

// NewPoWApp 创建 PoW 应用服务
func NewPoWApp(
	challengeRepo challenge.Repository,
	siteRepo site.Repository,
	ipPolicyRepo ippolicy.Repository,
	ipRecordRepo iprecord.Repository,
	verifylogRepo verifylog.Repository,
	log *infralogger.LogService,
) *PoWApp {
	return &PoWApp{
		challengeRepo: challengeRepo,
		siteRepo:      siteRepo,
		ipPolicyRepo:  ipPolicyRepo,
		ipRecordRepo:  ipRecordRepo,
		verifylogRepo: verifylogRepo,
		log:           log,
	}
}

// GetSiteForChallenge 根据 sitekey 获取站点并验证域名
func (app *PoWApp) GetSiteForChallenge(sitekey, origin string) (*site.Site, error) {
	s, err := app.siteRepo.FindBySitekey(sitekey)
	if err != nil {
		app.log.Warnf("[POW] GetSiteForChallenge: site not found by sitekey=%s", sitekey)
		return nil, ErrSiteNotFound
	}
	if !s.Enabled {
		app.log.Warnf("[POW] GetSiteForChallenge: site disabled (siteID=%d, sitekey=%s)", s.ID, sitekey)
		return nil, ErrSiteDisabled
	}
	if !s.ValidateDomain(origin) {
		app.log.Warnf("[POW] GetSiteForChallenge: domain mismatch (siteID=%d, origin=%s)", s.ID, origin)
		return nil, ErrDomainMismatch
	}
	app.log.Infof("[POW] GetSiteForChallenge: success (siteID=%d, sitekey=%s, origin=%s)", s.ID, sitekey, origin)
	return s, nil
}

// CreateChallenge 创建单个 PoW challenge
func (app *PoWApp) CreateChallenge(
	s *site.Site,
	clientIP,
	origin,
	action string,
	difficulty int,
) (*challenge.Challenge, error) {

	challengeStr, err := generateRandomString(32)
	if err != nil {
		app.log.Errorf("[POW] CreateChallenge: failed to generate random string: %v", err)
		return nil, err
	}

	c := challenge.NewChallenge(s.ID, challengeStr, difficulty, clientIP, origin, action,
		time.Now().Add(time.Duration(s.ChallengeTimeout)*time.Second))

	if err := app.challengeRepo.Create(c); err != nil {
		app.log.Errorf("[POW] CreateChallenge: DB error (siteID=%d, action=%s): %v", s.ID, action, err)
		return nil, err
	}

	// 增加请求计数
	_ = app.incrementRequestCount(s.ID, clientIP)

	app.log.Infof("[POW] Challenge created (siteID=%d, challengeID=%s, difficulty=%d, action=%s)", s.ID, c.ID, difficulty, action)
	return c, nil
}

// CreateBatchChallenges 批量创建 challenges
func (app *PoWApp) CreateBatchChallenges(
	s *site.Site,
	clientIP,
	origin,
	action string,
	difficulty,
	batchCount int,
) ([]*challenge.Challenge, uuid.UUID, error) {

	batchID := uuid.New()
	challenges := make([]*challenge.Challenge, batchCount)
	expiresAt := time.Now().Add(time.Duration(s.ChallengeTimeout) * time.Second)

	for i := 0; i < batchCount; i++ {
		challengeStr, err := generateRandomString(32)
		if err != nil {
			app.log.Errorf("[POW] CreateBatchChallenges: failed to generate random string: %v", err)
			return nil, uuid.Nil, err
		}

		c := &challenge.Challenge{
			ID:         uuid.New(),
			SiteID:     s.ID,
			Challenge:  challengeStr,
			Difficulty: difficulty,
			ClientIP:   clientIP,
			Origin:     origin,
			Action:     action,
			ExpiresAt:  expiresAt,
			BatchID:    &batchID,
			BatchIndex: i,
			BatchTotal: batchCount,
		}

		if err := app.challengeRepo.Create(c); err != nil {
			app.log.Errorf("[POW] CreateBatchChallenges: DB error (siteID=%d, batchIndex=%d): %v", s.ID, i, err)
			return nil, uuid.Nil, err
		}

		challenges[i] = c
	}

	_ = app.incrementRequestCount(s.ID, clientIP)

	app.log.Infof("[POW] Batch challenges created (siteID=%d, batchID=%s, count=%d, difficulty=%d)", s.ID, batchID, batchCount, difficulty)
	return challenges, batchID, nil
}

// VerifyChallenge 验证并解决单个 challenge
func (app *PoWApp) VerifyChallenge(
	challengeID uuid.UUID,
	nonce,
	clientIP string,
	iterations int64,
	solveTimeMs int,
) (*challenge.Challenge, error) {

	c, err := app.challengeRepo.FindByIDWithSite(challengeID)
	if err != nil {
		app.log.Warnf("[POW] VerifyChallenge: challenge not found (challengeID=%s, clientIP=%s)", challengeID, clientIP)
		return nil, ErrChallengeNotFound
	}

	if c.IsExpired() {
		app.log.Warnf("[POW] VerifyChallenge: challenge expired (challengeID=%s, clientIP=%s, siteID=%d)", challengeID, clientIP, c.SiteID)
		app.logVerification(c, clientIP, iterations, solveTimeMs, false)
		_ = app.incrementFailureCount(c.SiteID, clientIP)
		return nil, ErrChallengeExpired
	}

	if c.Solved {
		app.log.Warnf("[POW] VerifyChallenge: challenge already solved (challengeID=%s, clientIP=%s)", challengeID, clientIP)
		return nil, ErrChallengeSolved
	}

	if !c.VerifySolution(nonce, c.Difficulty) {
		app.log.Warnf("[POW] VerifyChallenge: invalid solution (challengeID=%s, clientIP=%s, difficulty=%d)", challengeID, clientIP, c.Difficulty)
		app.logVerification(c, clientIP, iterations, solveTimeMs, false)
		_ = app.incrementFailureCount(c.SiteID, clientIP)
		return nil, ErrInvalidSolution
	}

	c.MarkSolved()
	if err := app.challengeRepo.Save(c); err != nil {
		app.log.Errorf("[POW] VerifyChallenge: save error (challengeID=%s): %v", challengeID, err)
		return nil, err
	}

	app.logVerification(c, clientIP, iterations, solveTimeMs, true)
	app.log.Infof("[POW] Challenge verified successfully (challengeID=%s, clientIP=%s, iterations=%d, solveTimeMs=%d)", challengeID, clientIP, iterations, solveTimeMs)

	return c, nil
}

// VerifyBatchChallenges 验证批量 challenges
func (app *PoWApp) VerifyBatchChallenges(
	batchID uuid.UUID,
	nonces []string,
	clientIP string,
	iterations int64,
	solveTimeMs int,
) (*challenge.Challenge, error) {

	challenges, err := app.challengeRepo.FindByBatchID(batchID)
	if err != nil {
		app.log.Warnf("[POW] VerifyBatchChallenges: batch not found (batchID=%s, clientIP=%s)", batchID, clientIP)
		return nil, ErrChallengeNotFound
	}

	if len(challenges) != len(nonces) {
		app.log.Warnf("[POW] VerifyBatchChallenges: nonce count mismatch (batchID=%s, expected=%d, got=%d)", batchID, len(challenges), len(nonces))
		return nil, ErrNoncesMismatch
	}

	var firstSolved *challenge.Challenge
	for i, c := range challenges {
		if c.IsExpired() {
			app.log.Warnf("[POW] VerifyBatchChallenges: challenge expired (batchID=%s, index=%d, challengeID=%s)", batchID, i, c.ID)
			app.logVerification(c, clientIP, iterations, solveTimeMs, false)
			_ = app.incrementFailureCount(c.SiteID, clientIP)
			continue
		}

		if c.Solved {
			continue
		}

		if !c.VerifySolution(nonces[i], c.Difficulty) {
			app.log.Warnf("[POW] VerifyBatchChallenges: invalid solution (batchID=%s, index=%d, challengeID=%s)", batchID, i, c.ID)
			app.logVerification(c, clientIP, iterations, solveTimeMs, false)
			_ = app.incrementFailureCount(c.SiteID, clientIP)
			continue
		}

		c.MarkSolved()
		if err := app.challengeRepo.Save(c); err != nil {
			app.log.Errorf("[POW] VerifyBatchChallenges: save error (batchID=%s, index=%d): %v", batchID, i, err)
			return nil, err
		}

		app.logVerification(c, clientIP, iterations, solveTimeMs, true)

		if firstSolved == nil {
			firstSolved = c
		}
	}

	if firstSolved == nil {
		app.log.Warnf("[POW] VerifyBatchChallenges: no challenge solved successfully (batchID=%s, clientIP=%s)", batchID, clientIP)
		return nil, ErrInvalidSolution
	}

	app.log.Infof("[POW] Batch verified successfully (batchID=%s, clientIP=%s, solvedChallenges=%d)", batchID, clientIP, len(challenges))
	return firstSolved, nil
}

// SiteVerifyToken 站点验证 token
func (app *PoWApp) SiteVerifyToken(secret, token string) (*challenge.Challenge, error) {
	s, err := app.siteRepo.FindBySecret(secret)
	if err != nil {
		app.log.Warnf("[POW] SiteVerifyToken: site not found by secret")
		return nil, ErrSiteNotFound
	}

	c, err := app.challengeRepo.FindBySiteIDAndToken(s.ID, token)
	if err != nil {
		app.log.Warnf("[POW] SiteVerifyToken: token not found (siteID=%d, token=%s)", s.ID, token[:8]+"...")
		return nil, ErrTokenNotFound
	}

	if c.TokenUsed {
		app.log.Warnf("[POW] SiteVerifyToken: token already used (siteID=%d, token=%s)", s.ID, token[:8]+"...")
		return nil, ErrTokenUsed
	}

	if c.IsExpired() {
		app.log.Warnf("[POW] SiteVerifyToken: token expired (siteID=%d, token=%s)", s.ID, token[:8]+"...")
		return nil, ErrTokenExpired
	}

	c.TokenUsed = true
	if err := app.challengeRepo.Save(c); err != nil {
		app.log.Errorf("[POW] SiteVerifyToken: save error (siteID=%d): %v", s.ID, err)
		return nil, err
	}

	app.log.Infof("[POW] Token verified successfully (siteID=%d, token=%s)", s.ID, token[:8]+"...")
	return c, nil
}

// CalculateDifficulty 计算动态难度
func (app *PoWApp) CalculateDifficulty(s *site.Site, clientIP string) int {
	policy, err := app.ipPolicyRepo.FindBySiteID(s.ID)
	if err != nil {
		app.log.Infof("[POW] CalculateDifficulty: no policy, using base (siteID=%d, baseDifficulty=%d)", s.ID, s.BaseDifficulty)
		return s.BaseDifficulty
	}

	record, err := app.getOrCreateIpRecord(s.ID, clientIP, policy)
	if err != nil {
		app.log.Infof("[POW] CalculateDifficulty: no IP record, using base (siteID=%d, baseDifficulty=%d)", s.ID, s.BaseDifficulty)
		return s.BaseDifficulty
	}

	difficulty := policy.GetEffectiveDifficulty(
		s.BaseDifficulty,
		s.MaxDifficulty,
		record.RequestCount,
		record.FailureCount,
	)

	app.log.Infof("[POW] Difficulty calculated (siteID=%d, clientIP=%s, difficulty=%d, requests=%d, failures=%d)", s.ID, clientIP, difficulty, record.RequestCount, record.FailureCount)
	return difficulty
}

// CalculateBatchCount 计算批量数量
func (app *PoWApp) CalculateBatchCount(s *site.Site, clientIP string) int {
	policy, err := app.ipPolicyRepo.FindBySiteID(s.ID)
	if err != nil {
		app.log.Infof("[POW] CalculateBatchCount: no policy, using min (siteID=%d, minBatchCount=%d)", s.ID, s.MinBatchCount)
		return s.MinBatchCount
	}

	record, err := app.getOrCreateIpRecord(s.ID, clientIP, policy)
	if err != nil {
		app.log.Infof("[POW] CalculateBatchCount: no IP record, using min (siteID=%d, minBatchCount=%d)", s.ID, s.MinBatchCount)
		return s.MinBatchCount
	}

	count := policy.GetEffectiveBatchCount(
		s.MinBatchCount,
		s.MaxBatchCount,
		record.RequestCount,
		record.FailureCount,
	)

	app.log.Infof("[POW] Batch count calculated (siteID=%d, clientIP=%s, count=%d, requests=%d, failures=%d)", s.ID, clientIP, count, record.RequestCount, record.FailureCount)
	return count
}

// getOrCreateIpRecord 获取或创建 IP 记录
func (app *PoWApp) getOrCreateIpRecord(
	siteID uint,
	clientIP string,
	policy *ippolicy.IpPolicy,
) (*iprecord.IpRecord, error) {

	record, err := app.ipRecordRepo.FindBySiteIDAndIP(siteID, clientIP)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			record = iprecord.NewIpRecord(siteID, clientIP)
			if err := app.ipRecordRepo.Create(record); err != nil {
				return nil, err
			}
			return record, nil
		}
		return nil, err
	}

	if record.IsWindowExpired(policy.TimeWindowSeconds) {
		record.ResetWindow()
		if err := app.ipRecordRepo.Save(record); err != nil {
			return nil, err
		}
	}

	return record, nil
}

// incrementRequestCount 增加请求计数
func (app *PoWApp) incrementRequestCount(siteID uint, clientIP string) error {
	policy, err := app.ipPolicyRepo.FindBySiteID(siteID)
	if err != nil {
		return nil // 不阻塞流程
	}

	record, err := app.getOrCreateIpRecord(siteID, clientIP, policy)
	if err != nil {
		return nil
	}

	record.IncrementRequest()
	return app.ipRecordRepo.Save(record)
}

// incrementFailureCount 增加失败计数
func (app *PoWApp) incrementFailureCount(siteID uint, clientIP string) error {
	policy, err := app.ipPolicyRepo.FindBySiteID(siteID)
	if err != nil {
		return nil
	}

	record, err := app.getOrCreateIpRecord(siteID, clientIP, policy)
	if err != nil {
		return nil
	}

	record.IncrementFailure()
	return app.ipRecordRepo.Save(record)
}

// logVerification 记录验证日志
func (app *PoWApp) logVerification(
	c *challenge.Challenge,
	clientIP string,
	iterations int64,
	solveTimeMs int,
	success bool,
) {
	log := verifylog.NewVerificationLog(c, clientIP, iterations, solveTimeMs, success)
	_ = app.verifylogRepo.Create(log)
}

// GetSiteStats 获取站点统计
func (app *PoWApp) GetSiteStats(siteID uint) (map[string]any, error) {
	totalChallenges, _ := app.challengeRepo.CountBySiteID(siteID)
	solvedChallenges, _ := app.challengeRepo.CountSolvedBySiteID(siteID)
	totalVerifications, _ := app.verifylogRepo.CountBySiteID(siteID)
	successVerifications, _ := app.verifylogRepo.CountSuccessBySiteID(siteID)

	yesterday := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	last24hChallenges, _ := app.challengeRepo.CountBySiteIDSince(siteID, yesterday)
	last24hSolved, _ := app.challengeRepo.CountSolvedBySiteIDSince(siteID, yesterday)

	return map[string]any{
		"total_challenges":      totalChallenges,
		"solved_challenges":     solvedChallenges,
		"total_verifications":   totalVerifications,
		"success_verifications": successVerifications,
		"last_24h_challenges":   last24hChallenges,
		"last_24h_solved":       last24hSolved,
	}, nil
}
