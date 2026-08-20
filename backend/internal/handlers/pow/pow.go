package pow

import (
	"net/http"

	"github.com/ctipscn/weauth/backend/internal/models"
	"github.com/ctipscn/weauth/backend/internal/services/ipstrategy"
	svcPow "github.com/ctipscn/weauth/backend/internal/services/pow"
	"github.com/ctipscn/weauth/backend/internal/services/website"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PowController PoW 控制器，封装 challenge 创建与验证接口。
// 所有依赖通过构造函数注入，不依赖全局变量。
type PowController struct {
	siteSvc *website.SiteService
	ipSvc   *ipstrategy.IPStrategyService
	powSvc  *svcPow.PowService
}

// NewPowController 创建 PowController 实例。
func NewPowController(siteSvc *website.SiteService, ipSvc *ipstrategy.IPStrategyService, powSvc *svcPow.PowService) *PowController {
	return &PowController{
		siteSvc: siteSvc,
		ipSvc:   ipSvc,
		powSvc:  powSvc,
	}
}

// CreateChallenge 创建 challenge
func (c *PowController) CreateChallenge(ctx *gin.Context) {
	var req CreateChallengeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	site, err := c.siteSvc.GetSiteBySitekey(req.Sitekey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !c.siteSvc.ValidateDomain(site, req.Origin) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Domain not allowed"})
		return
	}

	clientIP := ctx.ClientIP()

	if site.BatchModeEnabled {
		difficulty := site.BatchDifficulty
		batchCount, _ := c.ipSvc.CalculateBatchCount(site, clientIP)

		challenges, batchID, err := c.powSvc.CreateBatchChallenges(site, clientIP, req.Origin, req.Action, difficulty, batchCount)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create challenges"})
			return
		}

		challengeStrs := make([]string, len(challenges))
		for i, ch := range challenges {
			challengeStrs[i] = ch.Challenge
		}

		ctx.JSON(http.StatusOK, CreateChallengeResponse{
			ChallengeID: batchID.String(),
			Challenges:  challengeStrs,
			Difficulty:  difficulty,
			BatchMode:   true,
			BatchTotal:  batchCount,
		})
		return
	}

	difficulty, _ := c.ipSvc.CalculateDifficulty(site, clientIP)

	if req.Difficulty != nil && *req.Difficulty >= 1 && *req.Difficulty <= 10 {
		if len(req.Sitekey) >= 5 && req.Sitekey[:5] == "test-" {
			difficulty = *req.Difficulty
		}
	}

	challenge, err := c.powSvc.CreateChallenge(site, clientIP, req.Origin, req.Action, difficulty)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create challenge"})
		return
	}

	ctx.JSON(http.StatusOK, CreateChallengeResponse{
		ChallengeID: challenge.ID.String(),
		Challenge:   challenge.Challenge,
		Difficulty:  challenge.Difficulty,
		BatchMode:   false,
		BatchTotal:  1,
	})
}

// VerifyChallenge 验证 challenge
func (c *PowController) VerifyChallenge(ctx *gin.Context) {
	var req VerifyChallengeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	challengeID, err := uuid.Parse(req.ChallengeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid challenge_id"})
		return
	}

	clientIP := ctx.ClientIP()

	var challenge *models.Challenge

	if req.BatchMode {
		if len(req.Nonces) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "nonces required for batch mode"})
			return
		}

		var verifyErr error
		challenge, verifyErr = c.powSvc.VerifyBatchChallenges(challengeID, req.Nonces, clientIP, req.Iterations, req.SolveTimeMs)
		if verifyErr != nil {
			status := http.StatusBadRequest
			if verifyErr == svcPow.ErrChallengeNotFound {
				status = http.StatusNotFound
			}
			ctx.JSON(status, gin.H{"error": verifyErr.Error()})
			return
		}
	} else {
		if req.Nonce == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "nonce required"})
			return
		}

		var verifyErr error
		challenge, verifyErr = c.powSvc.VerifyChallenge(challengeID, req.Nonce, clientIP, req.Iterations, req.SolveTimeMs)
		if verifyErr != nil {
			status := http.StatusBadRequest
			if verifyErr == svcPow.ErrChallengeNotFound {
				status = http.StatusNotFound
			}
			ctx.JSON(status, gin.H{"error": verifyErr.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, VerifyChallengeResponse{
		Token: challenge.Token,
	})
}