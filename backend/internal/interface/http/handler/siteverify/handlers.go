package siteverify

import (
	"errors"
	"net/http"

	"github.com/ctipscn/weauth/backend/internal/application/pow"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"github.com/ctipscn/weauth/backend/internal/interface/http/request"
	"github.com/ctipscn/weauth/backend/internal/interface/http/response"

	"github.com/gin-gonic/gin"
)

// SiteVerifyHandler 站点验证 HTTP 处理器
type SiteVerifyHandler struct {
	powApp *pow.PoWApp
	log    *infralogger.LogService
}

// NewSiteVerifyHandler 创建站点验证 HTTP 处理器
func NewSiteVerifyHandler(powApp *pow.PoWApp, log *infralogger.LogService) *SiteVerifyHandler {
	return &SiteVerifyHandler{powApp: powApp, log: log}
}

// SiteVerify 站点验证 token
func (h *SiteVerifyHandler) SiteVerify(c *gin.Context) {
	_, req, err := request.BindJSON[SiteVerifyRequest](c)
	if err != nil {
		h.log.Warnf("[HTTP-VERIFY] SiteVerify: invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "invalid-input-response"))
		return
	}

	h.log.Infof(
		"[HTTP-VERIFY] SiteVerify entry (secret=%s, token=%s)", 
		safeTruncate(req.Secret, 4), 
		safeTruncate(req.Token, 8),
	)

	challenge, err := h.powApp.SiteVerifyToken(req.Secret, req.Token)
	if err != nil {
		errorCode := "invalid-input-response"
		switch {
		case errors.Is(err, pow.ErrSiteNotFound), errors.Is(err, pow.ErrInvalidSecret):
			errorCode = "invalid-input-secret"
		case errors.Is(err, pow.ErrTokenNotFound):
			errorCode = "invalid-input-response"
		case errors.Is(err, pow.ErrTokenUsed):
			errorCode = "timeout-or-duplicate"
		case errors.Is(err, pow.ErrTokenExpired):
			errorCode = "timeout-or-duplicate"
		}

		h.log.Warnf("[HTTP-VERIFY] SiteVerify failed: %s (errorCode=%s)", err.Error(), errorCode)
		c.JSON(http.StatusOK, response.ErrorBaseResponse(0, 0, errorCode))
		return
	}

	h.log.Infof(
		"[HTTP-VERIFY] SiteVerify success (hostname=%s, action=%s)", 
		challenge.Origin, 
		challenge.Action,
	)

	c.JSON(http.StatusOK, response.SuccessBaseResponse(gin.H{
		"hostname":     challenge.Origin,
		"action":       challenge.Action,
		"challenge_ts": challenge.CreatedAt.Format("2006-01-02T15:04:05Z"),
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
