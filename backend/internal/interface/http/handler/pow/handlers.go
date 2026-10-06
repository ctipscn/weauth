package pow

import (
	"errors"
	"net/http"

	"github.com/ctipscn/weauth/backend/internal/application/pow"
	"github.com/ctipscn/weauth/backend/internal/domain/challenge"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"github.com/ctipscn/weauth/backend/internal/interface/http/request"
	"github.com/ctipscn/weauth/backend/internal/interface/http/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PoWHandler PoW HTTP 处理器
type PoWHandler struct {
	powApp *pow.PoWApp
	log    *infralogger.LogService
}

// NewPoWHandler 创建 PoW HTTP 处理器
func NewPoWHandler(powApp *pow.PoWApp, log *infralogger.LogService) *PoWHandler {
	return &PoWHandler{powApp: powApp, log: log}
}

// CreateChallenge 创建 challenge
func (h *PoWHandler) CreateChallenge(c *gin.Context) {
	_, req, err := request.BindJSON[CreateChallengeRequest](c)
	if err != nil {
		h.log.Warnf("[HTTP-POW] CreateChallenge: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid request: "+err.Error()))
		return
	}

	h.log.Infof(
		"[HTTP-POW] CreateChallenge (sitekey=%s, origin=%s)", 
		safeTruncate(req.Sitekey, 8), 
		req.Origin,
	)

	// 获取站点并验证域名（应用层处理）
	site, err := h.powApp.GetSiteForChallenge(req.Sitekey, req.Origin)
	if err != nil {
		h.log.Warnf(
			"[HTTP-POW] CreateChallenge: site validation failed: %v (sitekey=%s, origin=%s)", 
			err, 
			safeTruncate(req.Sitekey, 8), 
			req.Origin,
		)
		
		switch {
		case errors.Is(err, pow.ErrSiteNotFound):
			c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, err.Error()))
		case errors.Is(err, pow.ErrSiteDisabled):
			c.JSON(http.StatusForbidden, response.ErrorBaseResponse(0, 0, err.Error()))
		case errors.Is(err, pow.ErrDomainMismatch):
			c.JSON(http.StatusForbidden, response.ErrorBaseResponse(0, 0, err.Error()))
		default:
			c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, err.Error()))
		}
		return
	}

	// 获取客户端 IP
	clientIP := c.ClientIP()

	// 检查是否启用批量模式
	if site.BatchModeEnabled {
		difficulty := site.BatchDifficulty
		batchCount := h.powApp.CalculateBatchCount(site, clientIP)

		challenges, batchID, err := h.powApp.CreateBatchChallenges(
			site, clientIP, req.Origin, req.Action, difficulty, batchCount,
		)

		if err != nil {
			h.log.Errorf("[HTTP-POW] CreateChallenge batch failed: %v (sitekey=%s)", 
			err, 
			safeTruncate(req.Sitekey, 8),
		)
			c.JSON(
				http.StatusInternalServerError, 
				response.ErrorBaseResponse(0, 0, "Failed to create challenges"),
			)
			return
		}

		challengeStrs := make([]string, len(challenges))
		for i, ch := range challenges {
			challengeStrs[i] = ch.Challenge
		}

		h.log.Infof(
			"[HTTP-POW] CreateChallenge batch success (batchID=%s, count=%d)", 
			batchID.String(), 
			len(challenges),
		)
		c.JSON(http.StatusOK, response.SuccessBaseResponse(CreateChallengeResponse{
			ChallengeID: batchID.String(),
			Challenges:  challengeStrs,
			Difficulty:  difficulty,
			BatchMode:   true,
			BatchTotal:  batchCount,
		}))
		return
	}

	// 单个模式
	difficulty := h.powApp.CalculateDifficulty(site, clientIP)

	ch, err := h.powApp.CreateChallenge(site, clientIP, req.Origin, req.Action, difficulty)
	if err != nil {
		h.log.Errorf("[HTTP-POW] CreateChallenge failed: %v (sitekey=%s)", err, safeTruncate(req.Sitekey, 8))
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to create challenge"))
		return
	}

	h.log.Infof("[HTTP-POW] CreateChallenge success (challengeID=%s, difficulty=%d)", ch.ID.String(), ch.Difficulty)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(CreateChallengeResponse{
		ChallengeID: ch.ID.String(),
		Challenge:   ch.Challenge,
		Difficulty:  ch.Difficulty,
		BatchMode:   false,
		BatchTotal:  1,
	}))
}

// VerifyChallenge 验证 challenge
func (h *PoWHandler) VerifyChallenge(c *gin.Context) {
	_, req, err := request.BindJSON[VerifyChallengeRequest](c)
	if err != nil {
		h.log.Warnf("[HTTP-POW] VerifyChallenge: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid request: "+err.Error()))
		return
	}

	challengeID, err := uuid.Parse(req.ChallengeID)
	if err != nil {
		h.log.Warnf("[HTTP-POW] VerifyChallenge: invalid challengeID: %v", err)
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid challenge_id"))
		return
	}

	h.log.Infof("[HTTP-POW] VerifyChallenge (challengeID=%s, batchMode=%v)", challengeID.String(), req.BatchMode)
	clientIP := c.ClientIP()

	var verifiedChallenge *challenge.Challenge
	if req.BatchMode {
		if len(req.Nonces) == 0 {
			c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "nonces required for batch mode"))
			return
		}

		verifiedChallenge, err = h.powApp.VerifyBatchChallenges(challengeID, req.Nonces, clientIP, req.Iterations, req.SolveTimeMs)
	} else {
		if req.Nonce == "" {
			c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "nonce required"))
			return
		}

		verifiedChallenge, err = h.powApp.VerifyChallenge(challengeID, req.Nonce, clientIP, req.Iterations, req.SolveTimeMs)
	}

	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, pow.ErrChallengeNotFound) {
			status = http.StatusNotFound
		}
		h.log.Warnf("[HTTP-POW] VerifyChallenge failed (challengeID=%s, err=%v)", challengeID.String(), err)
		c.JSON(status, response.ErrorBaseResponse(0, 0, err.Error()))
		return
	}

	h.log.Infof("[HTTP-POW] VerifyChallenge success (challengeID=%s)", challengeID.String())
	c.JSON(http.StatusOK, response.SuccessBaseResponse(VerifyChallengeResponse{
		Token: verifiedChallenge.Token,
	}))
}

// safeTruncate 截断字符串用于日志，避免记录完整敏感信息
func safeTruncate(s string, maxLen int) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}