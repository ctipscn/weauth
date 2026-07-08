package handlers

import (
	"github.com/ctipscn/weauth/backend/models"
	"github.com/ctipscn/weauth/backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateChallengeRequest 创建 challenge 请求
type CreateChallengeRequest struct {
	Sitekey    string `json:"sitekey" binding:"required"`
	Origin     string `json:"origin" binding:"required"`
	Action     string `json:"action"`
	Difficulty *int   `json:"difficulty"` // 可选，用于测试指定难度
}

// CreateChallengeResponse 创建 challenge 响应
type CreateChallengeResponse struct {
	ChallengeID string   `json:"challenge_id"`
	Challenge   string   `json:"challenge,omitempty"`   // 单个模式使用
	Challenges  []string `json:"challenges,omitempty"`  // 批量模式使用
	Difficulty  int      `json:"difficulty"`
	BatchMode   bool     `json:"batch_mode"`            // 是否为批量模式
	BatchTotal  int      `json:"batch_total,omitempty"` // 批量模式的总数量
}

// CreateChallenge 创建 challenge
func CreateChallenge(c *gin.Context) {
	var req CreateChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 获取站点
	site, err := services.GetSiteBySitekey(req.Sitekey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证域名
	if !services.ValidateDomain(site, req.Origin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Domain not allowed"})
		return
	}

	// 获取客户端 IP
	clientIP := c.ClientIP()

	// 检查是否启用批量模式
	if site.BatchModeEnabled {
		// 批量模式：固定难度，动态数量
		difficulty := site.BatchDifficulty
		batchCount, _ := services.CalculateBatchCount(site, clientIP)

		// 创建批量 challenges
		challenges, batchID, err := services.CreateBatchChallenges(site, clientIP, req.Origin, req.Action, difficulty, batchCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create challenges"})
			return
		}

		// 收集所有 challenge 字符串
		challengeStrs := make([]string, len(challenges))
		for i, ch := range challenges {
			challengeStrs[i] = ch.Challenge
		}

		c.JSON(http.StatusOK, CreateChallengeResponse{
			ChallengeID: batchID.String(),
			Challenges:  challengeStrs,
			Difficulty:  difficulty,
			BatchMode:   true,
			BatchTotal:  batchCount,
		})
		return
	}

	// 单个模式（原有逻辑）
	// 计算动态难度
	difficulty, _ := services.CalculateDifficulty(site, clientIP)

	// 只允许测试站点使用自定义难度（安全限制）
	if req.Difficulty != nil && *req.Difficulty >= 1 && *req.Difficulty <= 10 {
		// 只有 sitekey 以 "test-" 开头的站点才允许自定义难度
		if len(req.Sitekey) >= 5 && req.Sitekey[:5] == "test-" {
			difficulty = *req.Difficulty
		}
	}

	// 创建 challenge
	challenge, err := services.CreateChallenge(site, clientIP, req.Origin, req.Action, difficulty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create challenge"})
		return
	}

	c.JSON(http.StatusOK, CreateChallengeResponse{
		ChallengeID: challenge.ID.String(),
		Challenge:   challenge.Challenge,
		Difficulty:  challenge.Difficulty,
		BatchMode:   false,
		BatchTotal:  1,
	})
}

// VerifyChallengeRequest 验证 challenge 请求
type VerifyChallengeRequest struct {
	ChallengeID string   `json:"challenge_id" binding:"required"`
	Nonce       string   `json:"nonce"`         // 单个模式使用
	Nonces      []string `json:"nonces"`        // 批量模式使用
	Iterations  int64    `json:"iterations"`
	SolveTimeMs int      `json:"solve_time_ms"`
	BatchMode   bool     `json:"batch_mode"`    // 是否为批量模式
}

// VerifyChallengeResponse 验证 challenge 响应
type VerifyChallengeResponse struct {
	Token string `json:"token"`
}

// VerifyChallenge 验证 challenge
func VerifyChallenge(c *gin.Context) {
	var req VerifyChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 解析 challenge ID (单个模式是 challenge ID，批量模式是 batch ID)
	challengeID, err := uuid.Parse(req.ChallengeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid challenge_id"})
		return
	}

	// 获取客户端 IP
	clientIP := c.ClientIP()

	var challenge *models.Challenge

	if req.BatchMode {
		// 批量模式验证
		if len(req.Nonces) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "nonces required for batch mode"})
			return
		}

		var verifyErr error
		challenge, verifyErr = services.VerifyBatchChallenges(challengeID, req.Nonces, clientIP, req.Iterations, req.SolveTimeMs)
		if verifyErr != nil {
			status := http.StatusBadRequest
			if verifyErr == services.ErrChallengeNotFound {
				status = http.StatusNotFound
			}
			c.JSON(status, gin.H{"error": verifyErr.Error()})
			return
		}
	} else {
		// 单个模式验证（原有逻辑）
		if req.Nonce == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "nonce required"})
			return
		}

		var verifyErr error
		challenge, verifyErr = services.VerifyChallenge(challengeID, req.Nonce, clientIP, req.Iterations, req.SolveTimeMs)
		if verifyErr != nil {
			status := http.StatusBadRequest
			if verifyErr == services.ErrChallengeNotFound {
				status = http.StatusNotFound
			}
			c.JSON(status, gin.H{"error": verifyErr.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, VerifyChallengeResponse{
		Token: challenge.Token,
	})
}
