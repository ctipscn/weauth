package pow

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/ctipscn/weauth/backend/internal/database"
	"github.com/ctipscn/weauth/backend/internal/models"
	"github.com/ctipscn/weauth/backend/internal/services/ipstrategy"
	"github.com/ctipscn/weauth/backend/internal/services/website"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PowService PoW 验证服务，封装 challenge 创建、验证与 token 管理。
type PowService struct {
	store   database.Storer
	siteSvc *website.SiteService
	ipSvc   *ipstrategy.IPStrategyService
}

// NewPowService 创建 PowService 实例。
func NewPowService(store database.Storer, siteSvc *website.SiteService, ipSvc *ipstrategy.IPStrategyService) *PowService {
	return &PowService{store: store, siteSvc: siteSvc, ipSvc: ipSvc}
}

// CreateChallenge 创建 challenge
func (s *PowService) CreateChallenge(site *models.Site, clientIP, origin, action string, difficulty int) (*models.Challenge, error) {
	db := s.store.DB()
	challengeStr, err := generateRandomString(32)
	if err != nil {
		return nil, err
	}

	challenge := &models.Challenge{
		SiteID:     site.ID,
		Challenge:  challengeStr,
		Difficulty: difficulty,
		ClientIP:   clientIP,
		Origin:     origin,
		Action:     action,
		ExpiresAt:  time.Now().Add(time.Duration(site.ChallengeTimeout) * time.Second),
	}

	if err := db.Create(challenge).Error; err != nil {
		return nil, err
	}

	_ = s.ipSvc.IncrementRequestCount(site.ID, clientIP)

	return challenge, nil
}

// CreateBatchChallenges 批量创建 challenges（用于批量模式）
func (s *PowService) CreateBatchChallenges(site *models.Site, clientIP, origin, action string, difficulty, batchCount int) ([]*models.Challenge, uuid.UUID, error) {
	db := s.store.DB()
	batchID := uuid.New()
	challenges := make([]*models.Challenge, batchCount)
	expiresAt := time.Now().Add(time.Duration(site.ChallengeTimeout) * time.Second)

	for i := 0; i < batchCount; i++ {
		challengeStr, err := generateRandomString(32)
		if err != nil {
			return nil, uuid.Nil, err
		}

		challenge := &models.Challenge{
			SiteID:     site.ID,
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

		if err := db.Create(challenge).Error; err != nil {
			return nil, uuid.Nil, err
		}

		challenges[i] = challenge
	}

	_ = s.ipSvc.IncrementRequestCount(site.ID, clientIP)

	return challenges, batchID, nil
}

// VerifyChallenge 验证 challenge 解决方案
func (s *PowService) VerifyChallenge(challengeID uuid.UUID, nonce string, clientIP string, iterations int64, solveTimeMs int) (*models.Challenge, error) {
	db := s.store.DB()
	var challenge models.Challenge
	err := db.Preload("Site").Where("id = ?", challengeID).First(&challenge).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChallengeNotFound
		}
		return nil, err
	}

	if challenge.IsExpired() {
		s.logVerification(&challenge, clientIP, iterations, solveTimeMs, false)
		_ = s.ipSvc.IncrementFailureCount(challenge.SiteID, clientIP)
		return nil, ErrChallengeExpired
	}

	if challenge.Solved {
		return nil, ErrChallengeSolved
	}

	if !verifySolution(challenge.Challenge, nonce, challenge.Difficulty) {
		s.logVerification(&challenge, clientIP, iterations, solveTimeMs, false)
		_ = s.ipSvc.IncrementFailureCount(challenge.SiteID, clientIP)
		return nil, ErrInvalidSolution
	}

	challenge.MarkSolved()
	if err := db.Save(&challenge).Error; err != nil {
		return nil, err
	}

	s.logVerification(&challenge, clientIP, iterations, solveTimeMs, true)

	return &challenge, nil
}

// VerifyBatchChallenges 批量验证 challenges（用于批量模式）
func (s *PowService) VerifyBatchChallenges(batchID uuid.UUID, nonces []string, clientIP string, totalIterations int64, solveTimeMs int) (*models.Challenge, error) {
	db := s.store.DB()
	var challenges []models.Challenge
	err := db.Preload("Site").Where("batch_id = ?", batchID).Order("batch_index").Find(&challenges).Error
	if err != nil {
		return nil, err
	}

	if len(challenges) == 0 {
		return nil, ErrChallengeNotFound
	}

	if len(nonces) != len(challenges) {
		return nil, errors.New("nonces count mismatch")
	}

	if challenges[0].IsExpired() {
		s.logVerification(&challenges[0], clientIP, totalIterations, solveTimeMs, false)
		_ = s.ipSvc.IncrementFailureCount(challenges[0].SiteID, clientIP)
		return nil, ErrChallengeExpired
	}

	for _, c := range challenges {
		if c.Solved {
			return nil, ErrChallengeSolved
		}
	}

	for i, challenge := range challenges {
		if !verifySolution(challenge.Challenge, nonces[i], challenge.Difficulty) {
			s.logVerification(&challenge, clientIP, totalIterations, solveTimeMs, false)
			_ = s.ipSvc.IncrementFailureCount(challenge.SiteID, clientIP)
			return nil, ErrInvalidSolution
		}
	}

	challenges[0].MarkSolved()
	if err := db.Save(&challenges[0]).Error; err != nil {
		return nil, err
	}

	for i := 1; i < len(challenges); i++ {
		now := time.Now()
		challenges[i].Solved = true
		challenges[i].SolvedAt = &now
		db.Save(&challenges[i])
	}

	s.logVerification(&challenges[0], clientIP, totalIterations, solveTimeMs, true)

	return &challenges[0], nil
}

// SiteVerifyToken 站点验证 token
func (s *PowService) SiteVerifyToken(secret, token string) (*models.Challenge, error) {
	db := s.store.DB()
	site, err := s.siteSvc.GetSiteBySecret(secret)
	if err != nil {
		return nil, err
	}

	var challenge models.Challenge
	err = db.Where("site_id = ? AND token = ?", site.ID, token).First(&challenge).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	if challenge.TokenUsed {
		return nil, ErrTokenUsed
	}

	tokenExpireAt := challenge.SolvedAt.Add(time.Duration(site.TokenTimeout) * time.Second)
	if time.Now().After(tokenExpireAt) {
		return nil, ErrTokenExpired
	}

	challenge.TokenUsed = true
	if err := db.Save(&challenge).Error; err != nil {
		return nil, err
	}

	return &challenge, nil
}

// GetChallengeByID 根据 ID 获取 challenge
func (s *PowService) GetChallengeByID(id uuid.UUID) (*models.Challenge, error) {
	db := s.store.DB()
	var challenge models.Challenge
	err := db.First(&challenge, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChallengeNotFound
		}
		return nil, err
	}
	return &challenge, nil
}

// GetSiteStats 获取站点统计（包含权限检查）
func (s *PowService) GetSiteStats(siteID uint, userID uint) (map[string]any, error) {
	db := s.store.DB()
	var site models.Site
	if err := db.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, website.ErrSiteNotFound
		}
		return nil, err
	}

	var totalChallenges int64
	var solvedChallenges int64
	var totalVerifications int64
	var successVerifications int64

	db.Model(&models.Challenge{}).Where("site_id = ?", siteID).Count(&totalChallenges)
	db.Model(&models.Challenge{}).Where("site_id = ? AND solved = true", siteID).Count(&solvedChallenges)
	db.Model(&models.VerificationLog{}).Where("site_id = ?", siteID).Count(&totalVerifications)
	db.Model(&models.VerificationLog{}).Where("site_id = ? AND success = true", siteID).Count(&successVerifications)

	yesterday := time.Now().Add(-24 * time.Hour)
	var last24hChallenges int64
	var last24hSolved int64
	db.Model(&models.Challenge{}).Where("site_id = ? AND created_at > ?", siteID, yesterday).Count(&last24hChallenges)
	db.Model(&models.Challenge{}).Where("site_id = ? AND solved = true AND created_at > ?", siteID, yesterday).Count(&last24hSolved)

	return map[string]any{
		"total_challenges":      totalChallenges,
		"solved_challenges":     solvedChallenges,
		"total_verifications":   totalVerifications,
		"success_verifications": successVerifications,
		"last_24h_challenges":   last24hChallenges,
		"last_24h_solved":       last24hSolved,
	}, nil
}

// --- private helpers ---

func verifySolution(challenge, nonce string, difficulty int) bool {
	data := challenge + nonce
	hash := sha256.Sum256([]byte(data))
	hashHex := hex.EncodeToString(hash[:])
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hashHex, prefix)
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func (s *PowService) logVerification(challenge *models.Challenge, clientIP string, iterations int64, solveTimeMs int, success bool) {
	db := s.store.DB()
	log := &models.VerificationLog{
		ChallengeID: &challenge.ID,
		SiteID:      &challenge.SiteID,
		ClientIP:    clientIP,
		Origin:      challenge.Origin,
		Action:      challenge.Action,
		Iterations:  iterations,
		SolveTimeMs: solveTimeMs,
		Success:     success,
	}
	db.Create(log)
}
