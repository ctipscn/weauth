package handlers

import (
	"github.com/ctipscn/weauth/backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListSites 获取站点列表
func ListSites(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sites, err := services.ListSites(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list sites"})
		return
	}

	c.JSON(http.StatusOK, sites)
}

// CreateSiteRequest 创建站点请求
type CreateSiteRequest struct {
	Name             string   `json:"name" binding:"required"`
	Domains          []string `json:"domains"`
	BatchModeEnabled *bool    `json:"batch_mode_enabled"` // 默认 true（批量模式）
	BatchDifficulty  *int     `json:"batch_difficulty"`   // 批量模式难度，默认 4
	MinBatchCount    *int     `json:"min_batch_count"`    // 最小批次数，默认 1
	MaxBatchCount    *int     `json:"max_batch_count"`    // 最大批次数，默认 10
}

// CreateSite 创建站点
func CreateSite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 构建配置选项
	opts := services.CreateSiteOptions{
		BatchModeEnabled: req.BatchModeEnabled,
		BatchDifficulty:  req.BatchDifficulty,
		MinBatchCount:    req.MinBatchCount,
		MaxBatchCount:    req.MaxBatchCount,
	}

	site, err := services.CreateSite(userID.(uint), req.Name, req.Domains, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create site"})
		return
	}

	// 返回时包含 secret
	c.JSON(http.StatusCreated, gin.H{
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
func GetSite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(c.Param("id"))

	site, err := services.GetSiteByID(id, userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	c.JSON(http.StatusOK, site)
}

// UpdateSiteRequest 更新站点请求
type UpdateSiteRequest struct {
	Name             *string `json:"name"`
	BaseDifficulty   *int    `json:"base_difficulty"`
	MaxDifficulty    *int    `json:"max_difficulty"`
	ChallengeTimeout *int    `json:"challenge_timeout"`
	TokenTimeout     *int    `json:"token_timeout"`
	Enabled          *bool   `json:"enabled"`
	// 批量模式配置
	BatchModeEnabled *bool `json:"batch_mode_enabled"`
	BatchDifficulty  *int  `json:"batch_difficulty"`
	MinBatchCount    *int  `json:"min_batch_count"`
	MaxBatchCount    *int  `json:"max_batch_count"`
}

// UpdateSite 更新站点
func UpdateSite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(c.Param("id"))

	var req UpdateSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	updates := make(map[string]interface{})
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
	// 批量模式配置
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

	site, err := services.UpdateSite(id, userID.(uint), updates)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	c.JSON(http.StatusOK, site)
}

// DeleteSite 删除站点
func DeleteSite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(c.Param("id"))

	if err := services.DeleteSite(id, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete site"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Site deleted"})
}

// ListDomains 获取域名列表
func ListDomains(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(c.Param("id"))

	site, err := services.GetSiteByID(id, userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	c.JSON(http.StatusOK, site.Domains)
}

// AddDomainRequest 添加域名请求
type AddDomainRequest struct {
	Domain string `json:"domain" binding:"required"`
}

// AddDomain 添加域名
func AddDomain(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(c.Param("id"))

	var req AddDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	domain, err := services.AddDomain(id, userID.(uint), req.Domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add domain"})
		return
	}

	c.JSON(http.StatusCreated, domain)
}

// DeleteDomain 删除域名
func DeleteDomain(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	siteID := parseUint(c.Param("id"))
	domainID := parseUint(c.Param("domainId"))

	if err := services.DeleteDomain(siteID, domainID, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete domain"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Domain deleted"})
}

// GetSiteSecret 获取站点 Secret（需要管理员权限）
func GetSiteSecret(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(c.Param("id"))

	site, err := services.GetSiteByID(id, userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"secret": site.Secret,
	})
}

// RegenerateSecret 重新生成 secret
func RegenerateSecret(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(c.Param("id"))

	secret, err := services.RegenerateSiteSecret(id, userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"secret": secret})
}

// GetIpPolicy 获取 IP 策略
func GetIpPolicy(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(c.Param("id"))

	policy, err := services.GetIpPolicy(id, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get IP policy"})
		return
	}

	c.JSON(http.StatusOK, policy)
}

// UpdateIpPolicyRequest 更新 IP 策略请求
type UpdateIpPolicyRequest struct {
	TimeWindowSeconds   *int     `json:"time_window_seconds"`
	RequestThreshold    *int     `json:"request_threshold"`
	DifficultyIncrement *int     `json:"difficulty_increment"` // 动态难度模式：难度增量
	BatchIncrement      *int     `json:"batch_increment"`      // 批量模式：批次增量
	FailureWeight       *float64 `json:"failure_weight"`
	Enabled             *bool    `json:"enabled"`
}

// UpdateIpPolicy 更新 IP 策略
func UpdateIpPolicy(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(c.Param("id"))

	var req UpdateIpPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	updates := make(map[string]interface{})
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

	policy, err := services.UpdateIpPolicy(id, userID.(uint), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update IP policy"})
		return
	}

	c.JSON(http.StatusOK, policy)
}

// ListIpRecords 获取 IP 记录列表
func ListIpRecords(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(c.Param("id"))

	records, err := services.ListIpRecords(id, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list IP records"})
		return
	}

	c.JSON(http.StatusOK, records)
}

// DeleteIpRecord 删除 IP 记录
func DeleteIpRecord(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	siteID := parseUint(c.Param("id"))
	clientIP := c.Param("ip")

	if err := services.DeleteIpRecord(siteID, clientIP, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete IP record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "IP record deleted"})
}

// GetSiteStats 获取站点统计
func GetSiteStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := parseUint(c.Param("id"))

	stats, err := services.GetSiteStats(id, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// parseUint 解析 uint
func parseUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 32)
	return uint(i)
}
