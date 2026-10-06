package admin

import (
	"net/http"
	"strconv"

	appadmin "github.com/ctipscn/weauth/backend/internal/application/admin"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"github.com/ctipscn/weauth/backend/internal/interface/http/request"
	"github.com/ctipscn/weauth/backend/internal/interface/http/response"

	"github.com/gin-gonic/gin"
)

// AdminHandler 管理员 HTTP 处理器
type AdminHandler struct {
	adminApp *appadmin.AdminApp
	log      *infralogger.LogService
}

// NewAdminHandler 创建管理员 HTTP 处理器
func NewAdminHandler(adminApp *appadmin.AdminApp, log *infralogger.LogService) *AdminHandler {
	return &AdminHandler{adminApp: adminApp, log: log}
}

// ListUsers 获取所有用户列表（管理员专用）
func (h *AdminHandler) ListUsers(c *gin.Context) {
	h.log.Infof("[HTTP-ADMIN] ListUsers entry")

	users, err := h.adminApp.ListUsers()
	if err != nil {
		h.log.Errorf("[HTTP-ADMIN] ListUsers failed: %v", err)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to list users"))
		return
	}

	h.log.Infof("[HTTP-ADMIN] ListUsers success (count=%d)", len(users))
	c.JSON(http.StatusOK, response.SuccessBaseResponse(users))
}

// DeleteUser 删除用户（管理员专用）
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-ADMIN] DeleteUser: invalid user ID (param=%s)", c.Param("id"))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid user ID"))
		return
	}

	_, req, err := request.BindJSON[DeleteUserRequest](c)
	if err != nil {
		// 不要求请求体必须存在，没有 body 则默认为 cascade_sites=false
		req = DeleteUserRequest{CascadeSites: false}
	}

	h.log.Infof("[HTTP-ADMIN] DeleteUser (userID=%d, cascadeSites=%v)", id, req.CascadeSites)

	if err := h.adminApp.DeleteUser(id, req.CascadeSites); err != nil {
		switch err {
		case appadmin.ErrUserNotFound:
			h.log.Warnf("[HTTP-ADMIN] DeleteUser: user not found (userID=%d)", id)
			c.JSON(http.StatusNotFound, response.ErrorBaseResponse(0, 0, "User not found"))
		default:
			h.log.Errorf("[HTTP-ADMIN] DeleteUser failed (userID=%d): %v", id, err)
			c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to delete user"))
		}
		return
	}

	h.log.Infof("[HTTP-ADMIN] DeleteUser success (userID=%d)", id)
	c.JSON(http.StatusOK, response.SuccessBaseResponse("User deleted"))
}

// ListSites 获取所有站点列表（管理员专用）
func (h *AdminHandler) ListSites(c *gin.Context) {
	h.log.Infof("[HTTP-ADMIN] ListSites entry")

	sites, err := h.adminApp.ListSites()
	if err != nil {
		h.log.Errorf("[HTTP-ADMIN] ListSites failed: %v", err)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to list sites"))
		return
	}

	h.log.Infof("[HTTP-ADMIN] ListSites success (count=%d)", len(sites))
	c.JSON(http.StatusOK, response.SuccessBaseResponse(sites))
}

// GetSite 获取指定站点详情（管理员专用，无归属检查）
func (h *AdminHandler) GetSite(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-ADMIN] GetSite: invalid site ID (param=%s)", c.Param("id"))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	h.log.Infof("[HTTP-ADMIN] GetSite (siteID=%d)", id)
	site, err := h.adminApp.GetSiteByID(id)
	if err != nil {
		switch err {
		case appadmin.ErrSiteNotFound:
			h.log.Warnf("[HTTP-ADMIN] GetSite: site not found (siteID=%d)", id)
			c.JSON(http.StatusNotFound, response.ErrorBaseResponse(0, 0, "Site not found"))
		default:
			h.log.Errorf("[HTTP-ADMIN] GetSite failed (siteID=%d): %v", id, err)
			c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to get site"))
		}
		return
	}

	h.log.Infof("[HTTP-ADMIN] GetSite success (siteID=%d)", id)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(site))
}

// DeleteSite 删除指定站点（管理员专用）
func (h *AdminHandler) DeleteSite(c *gin.Context) {
	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-ADMIN] DeleteSite: invalid site ID (param=%s)", c.Param("id"))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	h.log.Infof("[HTTP-ADMIN] DeleteSite (siteID=%d)", id)
	if err := h.adminApp.DeleteSite(id); err != nil {
		switch err {
		case appadmin.ErrSiteNotFound:
			h.log.Warnf("[HTTP-ADMIN] DeleteSite: site not found (siteID=%d)", id)
			c.JSON(http.StatusNotFound, response.ErrorBaseResponse(0, 0, "Site not found"))
		default:
			h.log.Errorf("[HTTP-ADMIN] DeleteSite failed (siteID=%d): %v", id, err)
			c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to delete site"))
		}
		return
	}

	h.log.Infof("[HTTP-ADMIN] DeleteSite success (siteID=%d)", id)
	c.JSON(http.StatusOK, response.SuccessBaseResponse("Site deleted"))
}

// DeleteDomain 删除指定站点下的域名（管理员专用）
func (h *AdminHandler) DeleteDomain(c *gin.Context) {
	siteID, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-ADMIN] DeleteDomain: invalid site ID (param=%s)", c.Param("id"))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	domainID, err := parseUint(c.Param("domainId"))
	if err != nil {
		h.log.Warnf("[HTTP-ADMIN] DeleteDomain: invalid domain ID (param=%s)", c.Param("domainId"))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid domain ID"))
		return
	}

	h.log.Infof("[HTTP-ADMIN] DeleteDomain (siteID=%d, domainID=%d)", siteID, domainID)
	if err := h.adminApp.DeleteDomain(siteID, domainID); err != nil {
		switch err {
		case appadmin.ErrSiteNotFound:
			h.log.Warnf("[HTTP-ADMIN] DeleteDomain: site not found (siteID=%d)", siteID)
			c.JSON(http.StatusNotFound, response.ErrorBaseResponse(0, 0, "Site not found"))
		case appadmin.ErrDomainNotFound:
			h.log.Warnf("[HTTP-ADMIN] DeleteDomain: domain not found (siteID=%d, domainID=%d)", siteID, domainID)
			c.JSON(http.StatusNotFound, response.ErrorBaseResponse(0, 0, "Domain not found"))
		default:
			h.log.Errorf("[HTTP-ADMIN] DeleteDomain failed (siteID=%d, domainID=%d): %v", siteID, domainID, err)
			c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to delete domain"))
		}
		return
	}

	h.log.Infof("[HTTP-ADMIN] DeleteDomain success (siteID=%d, domainID=%d)", siteID, domainID)
	c.JSON(http.StatusOK, response.SuccessBaseResponse("Domain deleted"))
}

// parseUint 字符串转 uint
func parseUint(s string) (uint, error) {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}