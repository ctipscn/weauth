package admin

import (
	"net/http"
	"strconv"

	"github.com/ctipscn/weauth/backend/internal/services/ipstrategy"
	"github.com/ctipscn/weauth/backend/internal/services/pow"
	"github.com/ctipscn/weauth/backend/internal/services/website"

	"github.com/gin-gonic/gin"
)

// AdminController 管理后台控制器，封装站点管理、域名管理、IP 策略与统计等接口。
// 所有依赖通过构造函数注入，不依赖全局变量。
type AdminController struct {
	siteSvc *website.SiteService
	ipSvc   *ipstrategy.IPStrategyService
	powSvc  *pow.PowService
}

// NewAdminController 创建 AdminController 实例。
func NewAdminController(siteSvc *website.SiteService, ipSvc *ipstrategy.IPStrategyService, powSvc *pow.PowService) *AdminController {
	return &AdminController{
		siteSvc: siteSvc,
		ipSvc:   ipSvc,
		powSvc:  powSvc,
	}
}

// ListSites 获取站点列表
func (c *AdminController) ListSites(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sites, err := c.siteSvc.ListSites(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list sites"})
		return
	}

	ctx.JSON(http.StatusOK, sites)
}

// CreateSite 创建站点
func (c *AdminController) CreateSite(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateSiteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	opts := website.CreateSiteOptions{
		BatchModeEnabled: req.BatchModeEnabled,
		BatchDifficulty:  req.BatchDifficulty,
		MinBatchCount:    req.MinBatchCount,
		MaxBatchCount:    req.MaxBatchCount,
	}

	site, err := c.siteSvc.CreateSite(userID.(uint), req.Name, req.Domains, opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create site"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
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
	})
}

// GetSite 获取站点详情
func (c *AdminController) GetSite(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(ctx.Param("id"))

	site, err := c.siteSvc.GetSiteByID(id, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	ctx.JSON(http.StatusOK, site)
}

// UpdateSite 更新站点
func (c *AdminController) UpdateSite(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(ctx.Param("id"))

	var req UpdateSiteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

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

	site, err := c.siteSvc.UpdateSite(id, userID.(uint), updates)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	ctx.JSON(http.StatusOK, site)
}

// DeleteSite 删除站点
func (c *AdminController) DeleteSite(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(ctx.Param("id"))

	if err := c.siteSvc.DeleteSite(id, userID.(uint)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete site"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Site deleted"})
}

// ListDomains 获取域名列表
func (c *AdminController) ListDomains(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(ctx.Param("id"))

	site, err := c.siteSvc.GetSiteByID(id, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	ctx.JSON(http.StatusOK, site.Domains)
}

// AddDomain 添加域名
func (c *AdminController) AddDomain(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(ctx.Param("id"))

	var req AddDomainRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	domain, err := c.siteSvc.AddDomain(id, userID.(uint), req.Domain)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add domain"})
		return
	}

	ctx.JSON(http.StatusCreated, domain)
}

// DeleteDomain 删除域名
func (c *AdminController) DeleteDomain(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	siteID := parseUint(ctx.Param("id"))
	domainID := parseUint(ctx.Param("domainId"))

	if err := c.siteSvc.DeleteDomain(siteID, domainID, userID.(uint)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete domain"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Domain deleted"})
}

// GetSiteSecret 获取站点 Secret
func (c *AdminController) GetSiteSecret(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(ctx.Param("id"))

	site, err := c.siteSvc.GetSiteByID(id, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"secret": site.Secret,
	})
}

// RegenerateSecret 重新生成 secret
func (c *AdminController) RegenerateSecret(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(ctx.Param("id"))

	secret, err := c.siteSvc.RegenerateSiteSecret(id, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"secret": secret})
}

// GetIpPolicy 获取 IP 策略
func (c *AdminController) GetIpPolicy(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(ctx.Param("id"))

	policy, err := c.ipSvc.GetIpPolicy(id, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get IP policy"})
		return
	}

	ctx.JSON(http.StatusOK, policy)
}

// UpdateIpPolicy 更新 IP 策略
func (c *AdminController) UpdateIpPolicy(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(ctx.Param("id"))

	var req UpdateIpPolicyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

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

	policy, err := c.ipSvc.UpdateIpPolicy(id, userID.(uint), updates)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update IP policy"})
		return
	}

	ctx.JSON(http.StatusOK, policy)
}

// ListIpRecords 获取 IP 记录列表
func (c *AdminController) ListIpRecords(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(ctx.Param("id"))

	records, err := c.ipSvc.ListIpRecords(id, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list IP records"})
		return
	}

	ctx.JSON(http.StatusOK, records)
}

// DeleteIpRecord 删除 IP 记录
func (c *AdminController) DeleteIpRecord(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	siteID := parseUint(ctx.Param("id"))
	clientIP := ctx.Param("ip")

	if err := c.ipSvc.DeleteIpRecord(siteID, clientIP, userID.(uint)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete IP record"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "IP record deleted"})
}

// GetSiteStats 获取站点统计
func (c *AdminController) GetSiteStats(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(ctx.Param("id"))

	stats, err := c.powSvc.GetSiteStats(id, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// parseUint 解析 uint
func parseUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 32)
	return uint(i)
}
