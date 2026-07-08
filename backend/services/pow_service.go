package services

import (
	"github.com/ctipscn/weauth/backend/database"
	"github.com/ctipscn/weauth/backend/models"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrChallengeNotFound = errors.New("challenge not found")
	ErrChallengeExpired  = errors.New("challenge expired")
	ErrChallengeSolved   = errors.New("challenge already solved")
	ErrInvalidSolution   = errors.New("invalid solution")
	ErrTokenNotFound     = errors.New("token not found")
	ErrTokenUsed         = errors.New("token already used")
	ErrTokenExpired      = errors.New("token expired")
)

// CreateChallenge 创建 challenge
func CreateChallenge(site *models.Site, clientIP, origin, action string, difficulty int) (*models.Challenge, error) {
	// 生成随机 challenge
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

	if err := database.DB.Create(challenge).Error; err != nil {
		return nil, err
	}

	// 增加请求计数
	_ = IncrementRequestCount(site.ID, clientIP)

	return challenge, nil
}

// CreateBatchChallenges 批量创建 challenges（用于批量模式）
func CreateBatchChallenges(site *models.Site, clientIP, origin, action string, difficulty, batchCount int) ([]*models.Challenge, uuid.UUID, error) {
	batchID := uuid.New()
	challenges := make([]*models.Challenge, batchCount)
	expiresAt := time.Now().Add(time.Duration(site.ChallengeTimeout) * time.Second)

	for i := 0; i < batchCount; i++ {
		// 生成随机 challenge
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

		if err := database.DB.Create(challenge).Error; err != nil {
			return nil, uuid.Nil, err
		}

		challenges[i] = challenge
	}

	// 增加请求计数
	_ = IncrementRequestCount(site.ID, clientIP)

	return challenges, batchID, nil
}

// VerifyChallenge 验证 challenge 解决方案
func VerifyChallenge(challengeID uuid.UUID, nonce string, clientIP string, iterations int64, solveTimeMs int) (*models.Challenge, error) {
	var challenge models.Challenge
	err := database.DB.Preload("Site").Where("id = ?", challengeID).First(&challenge).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChallengeNotFound
		}
		return nil, err
	}

	// 检查是否已过期
	if challenge.IsExpired() {
		// 记录失败
		logVerification(&challenge, clientIP, iterations, solveTimeMs, false)
		_ = IncrementFailureCount(challenge.SiteID, clientIP)
		return nil, ErrChallengeExpired
	}

	// 检查是否已解决
	if challenge.Solved {
		return nil, ErrChallengeSolved
	}

	// 验证解决方案
	if !verifySolution(challenge.Challenge, nonce, challenge.Difficulty) {
		// 记录失败
		logVerification(&challenge, clientIP, iterations, solveTimeMs, false)
		_ = IncrementFailureCount(challenge.SiteID, clientIP)
		return nil, ErrInvalidSolution
	}

	// 标记为已解决
	challenge.MarkSolved()
	if err := database.DB.Save(&challenge).Error; err != nil {
		return nil, err
	}

	// 记录成功
	logVerification(&challenge, clientIP, iterations, solveTimeMs, true)

	return &challenge, nil
}

// VerifyBatchChallenges 批量验证 challenges（用于批量模式）
func VerifyBatchChallenges(batchID uuid.UUID, nonces []string, clientIP string, totalIterations int64, solveTimeMs int) (*models.Challenge, error) {
	// 获取该批次的所有 challenges
	var challenges []models.Challenge
	err := database.DB.Preload("Site").Where("batch_id = ?", batchID).Order("batch_index").Find(&challenges).Error
	if err != nil {
		return nil, err
	}

	if len(challenges) == 0 {
		return nil, ErrChallengeNotFound
	}

	// 检查 nonces 数量是否匹配
	if len(nonces) != len(challenges) {
		return nil, errors.New("nonces count mismatch")
	}

	// 检查是否已过期（以第一个 challenge 为准）
	if challenges[0].IsExpired() {
		logVerification(&challenges[0], clientIP, totalIterations, solveTimeMs, false)
		_ = IncrementFailureCount(challenges[0].SiteID, clientIP)
		return nil, ErrChallengeExpired
	}

	// 检查是否已解决（任一已解决则整批已解决）
	for _, c := range challenges {
		if c.Solved {
			return nil, ErrChallengeSolved
		}
	}

	// 逐个验证解决方案
	for i, challenge := range challenges {
		if !verifySolution(challenge.Challenge, nonces[i], challenge.Difficulty) {
			// 记录失败
			logVerification(&challenge, clientIP, totalIterations, solveTimeMs, false)
			_ = IncrementFailureCount(challenge.SiteID, clientIP)
			return nil, ErrInvalidSolution
		}
	}

	// 所有验证通过，标记第一个 challenge 为已解决并生成 token
	challenges[0].MarkSolved()
	if err := database.DB.Save(&challenges[0]).Error; err != nil {
		return nil, err
	}

	// 标记其他 challenges 为已解决（不生成 token）
	for i := 1; i < len(challenges); i++ {
		now := time.Now()
		challenges[i].Solved = true
		challenges[i].SolvedAt = &now
		database.DB.Save(&challenges[i])
	}

	// 记录成功
	logVerification(&challenges[0], clientIP, totalIterations, solveTimeMs, true)

	return &challenges[0], nil
}

// verifySolution 验证解决方案
func verifySolution(challenge, nonce string, difficulty int) bool {
	data := challenge + nonce
	hash := sha256.Sum256([]byte(data))
	hashHex := hex.EncodeToString(hash[:])

	// 检查前导零
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hashHex, prefix)
}

// SiteVerifyToken 站点验证 token
func SiteVerifyToken(secret, token string) (*models.Challenge, error) {
	// 根据 secret 获取站点
	site, err := GetSiteBySecret(secret)
	if err != nil {
		return nil, err
	}

	// 查找 token 对应的 challenge
	var challenge models.Challenge
	err = database.DB.Where("site_id = ? AND token = ?", site.ID, token).First(&challenge).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	// 检查 token 是否已使用
	if challenge.TokenUsed {
		return nil, ErrTokenUsed
	}

	// 检查是否在 token 有效期内
	tokenExpireAt := challenge.SolvedAt.Add(time.Duration(site.TokenTimeout) * time.Second)
	if time.Now().After(tokenExpireAt) {
		return nil, ErrTokenExpired
	}

	// 标记 token 已使用
	challenge.TokenUsed = true
	if err := database.DB.Save(&challenge).Error; err != nil {
		return nil, err
	}

	return &challenge, nil
}

// GetChallengeByID 根据 ID 获取 challenge
func GetChallengeByID(id uuid.UUID) (*models.Challenge, error) {
	var challenge models.Challenge
	err := database.DB.First(&challenge, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChallengeNotFound
		}
		return nil, err
	}
	return &challenge, nil
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

// logVerification 记录验证日志
func logVerification(challenge *models.Challenge, clientIP string, iterations int64, solveTimeMs int, success bool) {
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
	database.DB.Create(log)
}

// GetSiteStats 获取站点统计（包含权限检查）
func GetSiteStats(siteID uint, userID uint) (map[string]interface{}, error) {
	// 检查站点是否属于该用户
	var site models.Site
	if err := database.DB.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSiteNotFound
		}
		return nil, err
	}

	var totalChallenges int64
	var solvedChallenges int64
	var totalVerifications int64
	var successVerifications int64

	database.DB.Model(&models.Challenge{}).Where("site_id = ?", siteID).Count(&totalChallenges)
	database.DB.Model(&models.Challenge{}).Where("site_id = ? AND solved = true", siteID).Count(&solvedChallenges)
	database.DB.Model(&models.VerificationLog{}).Where("site_id = ?", siteID).Count(&totalVerifications)
	database.DB.Model(&models.VerificationLog{}).Where("site_id = ? AND success = true", siteID).Count(&successVerifications)

	// 最近 24 小时
	yesterday := time.Now().Add(-24 * time.Hour)
	var last24hChallenges int64
	var last24hSolved int64
	database.DB.Model(&models.Challenge{}).Where("site_id = ? AND created_at > ?", siteID, yesterday).Count(&last24hChallenges)
	database.DB.Model(&models.Challenge{}).Where("site_id = ? AND solved = true AND created_at > ?", siteID, yesterday).Count(&last24hSolved)

	return map[string]interface{}{
		"total_challenges":      totalChallenges,
		"solved_challenges":     solvedChallenges,
		"total_verifications":   totalVerifications,
		"success_verifications": successVerifications,
		"last_24h_challenges":   last24hChallenges,
		"last_24h_solved":       last24hSolved,
	}, nil
}
