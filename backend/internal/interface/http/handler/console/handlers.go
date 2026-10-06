package console

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ctipscn/weauth/backend/internal/application/pow"
	"github.com/ctipscn/weauth/backend/internal/application/site"
	"github.com/ctipscn/weauth/backend/internal/domain"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"github.com/ctipscn/weauth/backend/internal/interface/http/request"
	"github.com/ctipscn/weauth/backend/internal/interface/http/response"

	"github.com/gin-gonic/gin"
)

// ConsoleHandler 站点管理 HTTP 处理器
type ConsoleHandler struct {
	siteApp *site.SiteApp
	powApp  *pow.PoWApp
	log     *infralogger.LogService
}

// NewConsoleHandler 创建站点管理 HTTP 处理器
func NewConsoleHandler(
	siteApp *site.SiteApp, 
	powApp *pow.PoWApp, 
	log *infralogger.LogService,
) *ConsoleHandler {
	return &ConsoleHandler{siteApp: siteApp, powApp: powApp, log: log}
}

// ListSites 获取站点列表
func (h *ConsoleHandler) ListSites(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] ListSites: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] ListSites: userID=%d", userID.(uint))
	sites, err := h.siteApp.ListSites(userID.(uint))
	if err != nil {
		h.log.Errorf("[HTTP-CONSOLE] ListSites failed: %v (userID=%d)", err, userID.(uint))
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to list sites"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] ListSites success (userID=%d, count=%d)", userID.(uint), len(sites))
	c.JSON(http.StatusOK, response.SuccessBaseResponse(sites))
}

// CreateSite 创建站点
func (h *ConsoleHandler) CreateSite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] CreateSite: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	_, req, err := request.BindJSON[CreateSiteRequest](c)
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] CreateSite: invalid request body (userID=%d, err=%v)", userID.(uint), err)
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid request"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] CreateSite entry (userID=%d, name=%s)", userID.(uint), req.Name)

	opts := site.CreateSiteOptions{
		BatchModeEnabled: req.BatchModeEnabled,
		BatchDifficulty:  req.BatchDifficulty,
		MinBatchCount:    req.MinBatchCount,
		MaxBatchCount:    req.MaxBatchCount,
	}

	site, err := h.siteApp.CreateSite(userID.(uint), req.Name, req.Domains, opts)
	if err != nil {
		h.log.Errorf("[HTTP-CONSOLE] CreateSite failed: %v (userID=%d, name=%s)", err, userID.(uint), req.Name)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to create site"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] CreateSite success (userID=%d, name=%s, siteID=%d)", userID.(uint), req.Name, site.ID)
	c.JSON(http.StatusCreated, response.SuccessBaseResponse(gin.H{
		"id":                  site.ID,
		"name":                site.Name,
		"sitekey":             site.Sitekey,
		"secret":              site.Secret,
		"base_difficulty":     site.BaseDifficulty,
		"max_difficulty":      site.MaxDifficulty,
		"challenge_timeout":   site.ChallengeTimeout,
		"token_timeout":       site.TokenTimeout,
		"enabled":             site.Enabled,
		"batch_mode_enabled":  site.BatchModeEnabled,
		"batch_difficulty":    site.BatchDifficulty,
		"min_batch_count":     site.MinBatchCount,
		"max_batch_count":     site.MaxBatchCount,
		"domains":             site.Domains,
		"created_at":          site.CreatedAt,
	}))
}

// GetSite 获取站点详情
func (h *ConsoleHandler) GetSite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] GetSite: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] GetSite: invalid site ID (userID=%d, param=%s)", userID.(uint), c.Param("id"))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] GetSite (userID=%d, siteID=%d)", userID.(uint), id)

	site, err := h.siteApp.GetSiteByID(id, userID.(uint))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] GetSite: site not found (userID=%d, siteID=%d, err=%v)", userID.(uint), id, err)
		c.JSON(http.StatusNotFound, response.ErrorBaseResponse(0, 0, "Site not found"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] GetSite success (userID=%d, siteID=%d, name=%s)", userID.(uint), id, site.Name)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(site))
}

// UpdateSite 更新站点
func (h *ConsoleHandler) UpdateSite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] UpdateSite: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] UpdateSite: invalid site ID (userID=%d)", userID.(uint))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	_, req, err := request.BindJSON[UpdateSiteRequest](c)
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] UpdateSite: invalid request body (userID=%d, siteID=%d)", userID.(uint), id)
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid request"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] UpdateSite (userID=%d, siteID=%d)", userID.(uint), id)

	updates := make(map[string]any)
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.BaseDifficulty != nil {
		updates["base_difficulty"] = *req.BaseDifficulty
	}
	if req.MaxDifficulty != nil {
		updates["max_difficulty"] = *req.MaxDifficulty
	}
	if req.ChallengeTimeout != nil {
		updates["challenge_timeout"] = *req.ChallengeTimeout
	}
	if req.TokenTimeout != nil {
		updates["token_timeout"] = *req.TokenTimeout
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if req.BatchModeEnabled != nil {
		updates["batch_mode_enabled"] = *req.BatchModeEnabled
	}
	if req.BatchDifficulty != nil {
		updates["batch_difficulty"] = *req.BatchDifficulty
	}
	if req.MinBatchCount != nil {
		updates["min_batch_count"] = *req.MinBatchCount
	}
	if req.MaxBatchCount != nil {
		updates["max_batch_count"] = *req.MaxBatchCount
	}

	site, err := h.siteApp.UpdateSite(id, userID.(uint), updates)
	if err != nil {
		h.log.Errorf("[HTTP-CONSOLE] UpdateSite failed: %v (userID=%d, siteID=%d)", err, userID.(uint), id)
		c.JSON(http.StatusNotFound, response.ErrorBaseResponse(0, 0, "Site not found"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] UpdateSite success (userID=%d, siteID=%d)", userID.(uint), id)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(site))
}

// DeleteSite 删除站点
func (h *ConsoleHandler) DeleteSite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] DeleteSite: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] DeleteSite: invalid site ID (userID=%d)", userID.(uint))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] DeleteSite (userID=%d, siteID=%d)", userID.(uint), id)
	if err := h.siteApp.DeleteSite(id, userID.(uint)); err != nil {
		h.log.Errorf("[HTTP-CONSOLE] DeleteSite failed: %v (userID=%d, siteID=%d)", err, userID.(uint), id)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to delete site"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] DeleteSite success (userID=%d, siteID=%d)", userID.(uint), id)
	c.JSON(http.StatusOK, response.SuccessBaseResponse("Site deleted"))
}

// ListDomains 获取域名列表
func (h *ConsoleHandler) ListDomains(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] ListDomains: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] ListDomains: invalid site ID (userID=%d)", userID.(uint))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] ListDomains (userID=%d, siteID=%d)", userID.(uint), id)
	site, err := h.siteApp.GetSiteByID(id, userID.(uint))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] ListDomains: site not found (userID=%d, siteID=%d)", userID.(uint), id)
		c.JSON(http.StatusNotFound, response.ErrorBaseResponse(0, 0, "Site not found"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] ListDomains success (userID=%d, siteID=%d, count=%d)", userID.(uint), id, len(site.Domains))
	c.JSON(http.StatusOK, response.SuccessBaseResponse(site.Domains))
}

// AddDomain 添加域名
func (h *ConsoleHandler) AddDomain(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] AddDomain: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] AddDomain: invalid site ID (userID=%d)", userID.(uint))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	_, req, err := request.BindJSON[AddDomainRequest](c)
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] AddDomain: invalid request body (userID=%d, siteID=%d)", userID.(uint), id)
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid request"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] AddDomain (userID=%d, siteID=%d, domain=%s)", userID.(uint), id, req.Domain)
	domainObj, err := h.siteApp.AddDomain(id, userID.(uint), req.Domain)
	if err != nil {
		h.log.Errorf("[HTTP-CONSOLE] AddDomain failed: %v (userID=%d, siteID=%d, domain=%s)", err, userID.(uint), id, req.Domain)
		if errors.Is(err, domain.ErrDomainAlreadyExists) {
			c.JSON(http.StatusConflict, response.ErrorBaseResponse(0, 0, "Domain already exists in whitelist"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to add domain"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] AddDomain success (userID=%d, siteID=%d, domainID=%d)", userID.(uint), id, domainObj.ID)
	c.JSON(http.StatusCreated, response.SuccessBaseResponse(domainObj))
}

// DeleteDomain 删除域名
func (h *ConsoleHandler) DeleteDomain(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] DeleteDomain: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	siteID, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] DeleteDomain: invalid site ID (userID=%d)", userID.(uint))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}
	domainID, err := parseUint(c.Param("domainId"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] DeleteDomain: invalid domain ID (userID=%d)", userID.(uint))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid domain ID"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] DeleteDomain (userID=%d, siteID=%d, domainID=%d)", userID.(uint), siteID, domainID)
	if err := h.siteApp.DeleteDomain(siteID, domainID, userID.(uint)); err != nil {
		h.log.Errorf("[HTTP-CONSOLE] DeleteDomain failed: %v (userID=%d, siteID=%d, domainID=%d)", err, userID.(uint), siteID, domainID)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to delete domain"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] DeleteDomain success (userID=%d, siteID=%d, domainID=%d)", userID.(uint), siteID, domainID)
	c.JSON(http.StatusOK, response.SuccessBaseResponse("Domain deleted"))
}

// GetSiteSecret 获取站点 Secret
func (h *ConsoleHandler) GetSiteSecret(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] GetSiteSecret: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] GetSiteSecret: invalid site ID (userID=%d)", userID.(uint))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] GetSiteSecret (userID=%d, siteID=%d)", userID.(uint), id)
	site, err := h.siteApp.GetSiteByID(id, userID.(uint))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] GetSiteSecret: site not found (userID=%d, siteID=%d)", userID.(uint), id)
		c.JSON(http.StatusNotFound, response.ErrorBaseResponse(0, 0, "Site not found"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] GetSiteSecret success (userID=%d, siteID=%d)", userID.(uint), id)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(gin.H{"secret": site.Secret}))
}

// RegenerateSecret 重新生成 secret
func (h *ConsoleHandler) RegenerateSecret(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] RegenerateSecret: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] RegenerateSecret: invalid site ID (userID=%d)", userID.(uint))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] RegenerateSecret (userID=%d, siteID=%d)", userID.(uint), id)
	secret, err := h.siteApp.RegenerateSiteSecret(id, userID.(uint))
	if err != nil {
		h.log.Errorf("[HTTP-CONSOLE] RegenerateSecret failed: %v (userID=%d, siteID=%d)", err, userID.(uint), id)
		c.JSON(http.StatusNotFound, response.ErrorBaseResponse(0, 0, "Site not found"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] RegenerateSecret success (userID=%d, siteID=%d)", userID.(uint), id)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(gin.H{"secret": secret}))
}

// GetIpPolicy 获取 IP 策略
func (h *ConsoleHandler) GetIpPolicy(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] GetIpPolicy: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] GetIpPolicy: invalid site ID (userID=%d)", userID.(uint))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] GetIpPolicy (userID=%d, siteID=%d)", userID.(uint), id)
	policy, err := h.siteApp.GetIpPolicy(id, userID.(uint))
	if err != nil {
		h.log.Errorf("[HTTP-CONSOLE] GetIpPolicy failed: %v (userID=%d, siteID=%d)", err, userID.(uint), id)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to get IP policy"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] GetIpPolicy success (userID=%d, siteID=%d)", userID.(uint), id)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(policy))
}

// UpdateIpPolicy 更新 IP 策略
func (h *ConsoleHandler) UpdateIpPolicy(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] UpdateIpPolicy: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] UpdateIpPolicy: invalid site ID (userID=%d)", userID.(uint))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	_, req, err := request.BindJSON[UpdateIpPolicyRequest](c)
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] UpdateIpPolicy: invalid request body (userID=%d, siteID=%d)", userID.(uint), id)
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid request"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] UpdateIpPolicy (userID=%d, siteID=%d)", userID.(uint), id)

	updates := make(map[string]any)
	if req.TimeWindowSeconds != nil {
		updates["time_window_seconds"] = *req.TimeWindowSeconds
	}
	if req.RequestThreshold != nil {
		updates["request_threshold"] = *req.RequestThreshold
	}
	if req.DifficultyIncrement != nil {
		updates["difficulty_increment"] = *req.DifficultyIncrement
	}
	if req.BatchIncrement != nil {
		updates["batch_increment"] = *req.BatchIncrement
	}
	if req.FailureWeight != nil {
		updates["failure_weight"] = *req.FailureWeight
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	policy, err := h.siteApp.UpdateIpPolicy(id, userID.(uint), updates)
	if err != nil {
		h.log.Errorf("[HTTP-CONSOLE] UpdateIpPolicy failed: %v (userID=%d, siteID=%d)", err, userID.(uint), id)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to update IP policy"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] UpdateIpPolicy success (userID=%d, siteID=%d)", userID.(uint), id)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(policy))
}

// ListIpRecords 获取 IP 记录列表
func (h *ConsoleHandler) ListIpRecords(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] ListIpRecords: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] ListIpRecords: invalid site ID (userID=%d)", userID.(uint))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] ListIpRecords (userID=%d, siteID=%d)", userID.(uint), id)
	records, err := h.siteApp.ListIpRecords(id, userID.(uint))
	if err != nil {
		h.log.Errorf("[HTTP-CONSOLE] ListIpRecords failed: %v (userID=%d, siteID=%d)", err, userID.(uint), id)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to list IP records"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] ListIpRecords success (userID=%d, siteID=%d, count=%d)", userID.(uint), id, len(records))
	c.JSON(http.StatusOK, response.SuccessBaseResponse(records))
}

// DeleteIpRecord 删除 IP 记录
func (h *ConsoleHandler) DeleteIpRecord(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] DeleteIpRecord: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	siteID, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] DeleteIpRecord: invalid site ID (userID=%d)", userID.(uint))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}
	clientIP := c.Param("ip")

	h.log.Infof("[HTTP-CONSOLE] DeleteIpRecord (userID=%d, siteID=%d, ip=%s)", userID.(uint), siteID, clientIP)
	if err := h.siteApp.DeleteIpRecord(siteID, userID.(uint), clientIP); err != nil {
		h.log.Errorf("[HTTP-CONSOLE] DeleteIpRecord failed: %v (userID=%d, siteID=%d, ip=%s)", err, userID.(uint), siteID, clientIP)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to delete IP record"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] DeleteIpRecord success (userID=%d, siteID=%d, ip=%s)", userID.(uint), siteID, clientIP)
	c.JSON(http.StatusOK, response.SuccessBaseResponse("IP record deleted"))
}

// GetSiteStats 获取站点统计
func (h *ConsoleHandler) GetSiteStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		h.log.Warnf("[HTTP-CONSOLE] GetSiteStats: unauthorized access")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
		return
	}

	id, err := parseUint(c.Param("id"))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] GetSiteStats: invalid site ID (userID=%d)", userID.(uint))
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid site ID"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] GetSiteStats (userID=%d, siteID=%d)", userID.(uint), id)

	// 权限检查：确认站点属于该用户
	_, err = h.siteApp.GetSiteByID(id, userID.(uint))
	if err != nil {
		h.log.Warnf("[HTTP-CONSOLE] GetSiteStats: site not found (userID=%d, siteID=%d)", userID.(uint), id)
		c.JSON(http.StatusNotFound, response.ErrorBaseResponse(0, 0, "Site not found"))
		return
	}

	stats, err := h.powApp.GetSiteStats(id)
	if err != nil {
		h.log.Errorf("[HTTP-CONSOLE] GetSiteStats failed: %v (userID=%d, siteID=%d)", err, userID.(uint), id)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to get stats"))
		return
	}

	h.log.Infof("[HTTP-CONSOLE] GetSiteStats success (userID=%d, siteID=%d)", userID.(uint), id)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(stats))
}

func parseUint(s string) (uint, error) {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}