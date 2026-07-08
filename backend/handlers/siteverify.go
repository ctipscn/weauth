package handlers

import (
	"github.com/ctipscn/weauth/backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SiteVerifyRequest 站点验证请求
type SiteVerifyRequest struct {
	Secret string `json:"secret" binding:"required"`
	Token  string `json:"token" binding:"required"`
}

// SiteVerifyResponse 站点验证响应
type SiteVerifyResponse struct {
	Success     bool   `json:"success"`
	Hostname    string `json:"hostname,omitempty"`
	Action      string `json:"action,omitempty"`
	ChallengeTs string `json:"challenge_ts,omitempty"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
}

// SiteVerify 站点验证 token
func SiteVerify(c *gin.Context) {
	var req SiteVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, SiteVerifyResponse{
			Success:    false,
			ErrorCodes: []string{"invalid-input-response"},
		})
		return
	}

	challenge, err := services.SiteVerifyToken(req.Secret, req.Token)
	if err != nil {
		errorCode := "invalid-input-response"
		switch err {
		case services.ErrSiteNotFound, services.ErrInvalidSecret:
			errorCode = "invalid-input-secret"
		case services.ErrTokenNotFound:
			errorCode = "invalid-input-response"
		case services.ErrTokenUsed:
			errorCode = "timeout-or-duplicate"
		case services.ErrTokenExpired:
			errorCode = "timeout-or-duplicate"
		}

		c.JSON(http.StatusOK, SiteVerifyResponse{
			Success:    false,
			ErrorCodes: []string{errorCode},
		})
		return
	}

	c.JSON(http.StatusOK, SiteVerifyResponse{
		Success:     true,
		Hostname:    challenge.Origin,
		Action:      challenge.Action,
		ChallengeTs: challenge.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}
