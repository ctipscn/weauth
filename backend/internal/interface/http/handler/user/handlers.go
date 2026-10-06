package user

import (
	"errors"
	"net/http"

	appuser "github.com/ctipscn/weauth/backend/internal/application/user"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"github.com/ctipscn/weauth/backend/internal/interface/http/request"
	"github.com/ctipscn/weauth/backend/internal/interface/http/response"
	"github.com/ctipscn/weauth/backend/internal/interface/http/middleware"
	"github.com/gin-gonic/gin"
)

// UserHandler 用户 HTTP 处理器
type UserHandler struct {
	userApp *appuser.UserApp
	log     *infralogger.LogService
}

// NewUserHandler 创建用户 HTTP 处理器
func NewUserHandler(userApp *appuser.UserApp, log *infralogger.LogService) *UserHandler {
	return &UserHandler{userApp: userApp, log: log}
}

// GetProfile 获取用户资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		h.log.Warnf("[HTTP-USER] GetProfile: not authenticated")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Not authenticated"))
		return
	}

	h.log.Infof("[HTTP-USER] GetProfile (userID=%d, username=%s)", user.ID, user.Username)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(ProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Avatar:      user.Avatar,
		Role:        user.Role.String(),
		CreatedAt:   user.CreatedAt,
	}))
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		h.log.Warnf("[HTTP-USER] UpdateProfile: not authenticated")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Not authenticated"))
		return
	}

	_, req, err := request.BindJSON[UpdateProfileRequest](c)
	if err != nil {
		h.log.Warnf("[HTTP-USER] UpdateProfile: invalid request body (userID=%d, err=%v)", currentUser.ID, err)
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid request"))
		return
	}

	h.log.Infof("[HTTP-USER] UpdateProfile (userID=%d, username=%s)", currentUser.ID, req.Username)
	user, err := h.userApp.UpdateProfile(currentUser.ID, req.Username, req.Avatar)
	if err != nil {
		if errors.Is(err, appuser.ErrUsernameTaken) {
			h.log.Warnf("[HTTP-USER] UpdateProfile failed: username taken (userID=%d, username=%s)", currentUser.ID, req.Username)
			c.JSON(http.StatusConflict, response.ErrorBaseResponse(0, 0, "Username already taken"))
			return
		}
		h.log.Errorf("[HTTP-USER] UpdateProfile failed: %v (userID=%d)", err, currentUser.ID)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to update profile"))
		return
	}

	h.log.Infof("[HTTP-USER] UpdateProfile success (userID=%d)", user.ID)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(ProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Avatar:      user.Avatar,
		Role:        user.Role.String(),
		CreatedAt:   user.CreatedAt,
	}))
}

// DeactivateAccount 注销当前用户账号
func (h *UserHandler) DeactivateAccount(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		h.log.Warnf("[HTTP-USER] DeactivateAccount: not authenticated")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Not authenticated"))
		return
	}

	h.log.Infof("[HTTP-USER] DeactivateAccount (userID=%d, username=%s)", currentUser.ID, currentUser.Username)

	if err := h.userApp.DeactivateUser(currentUser.ID); err != nil {
		if errors.Is(err, appuser.ErrUserAlreadyDeactivated) {
			h.log.Warnf("[HTTP-USER] DeactivateAccount: already deactivated (userID=%d)", currentUser.ID)
			c.JSON(http.StatusConflict, response.ErrorBaseResponse(0, 0, "User already deactivated"))
			return
		}
		h.log.Errorf("[HTTP-USER] DeactivateAccount failed (userID=%d): %v", currentUser.ID, err)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to deactivate account"))
		return
	}

	h.log.Infof("[HTTP-USER] DeactivateAccount success (userID=%d)", currentUser.ID)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(map[string]string{
		"message": "Account deactivated successfully",
	}))
}