package siteverify

import (
	"net/http"

	svcPow "github.com/ctipscn/weauth/backend/internal/services/pow"
	svcWebsite "github.com/ctipscn/weauth/backend/internal/services/website"

	"github.com/gin-gonic/gin"
)

// SiteVerifyController 站点验证控制器，封装 /api/pow/siteverify 接口。
// 所有依赖通过构造函数注入，不依赖全局变量。
type SiteVerifyController struct {
	powSvc *svcPow.PowService
}

// NewSiteVerifyController 创建 SiteVerifyController 实例。
func NewSiteVerifyController(powSvc *svcPow.PowService) *SiteVerifyController {
	return &SiteVerifyController{powSvc: powSvc}
}

// SiteVerify 站点验证 token
func (c *SiteVerifyController) SiteVerify(ctx *gin.Context) {
	var req SiteVerifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, SiteVerifyResponse{
			Success:    false,
			ErrorCodes: []string{"invalid-input-response"},
		})
		return
	}

	challenge, err := c.powSvc.SiteVerifyToken(req.Secret, req.Token)
	if err != nil {
		errorCode := "invalid-input-response"
		switch err {
		case svcWebsite.ErrSiteNotFound, svcWebsite.ErrInvalidSecret:
			errorCode = "invalid-input-secret"
		case svcPow.ErrTokenNotFound:
			errorCode = "invalid-input-response"
		case svcPow.ErrTokenUsed:
			errorCode = "timeout-or-duplicate"
		case svcPow.ErrTokenExpired:
			errorCode = "timeout-or-duplicate"
		}

		ctx.JSON(http.StatusOK, SiteVerifyResponse{
			Success:    false,
			ErrorCodes: []string{errorCode},
		})
		return
	}

	ctx.JSON(http.StatusOK, SiteVerifyResponse{
		Success:     true,
		Hostname:    challenge.Origin,
		Action:      challenge.Action,
		ChallengeTs: challenge.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}