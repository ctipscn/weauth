package site

import (
	"strings"
	"time"

	"github.com/ctipscn/weauth/backend/internal/domain/user"
	"github.com/google/uuid"
)

// Site 站点实体
type Site struct {
	ID               uint      `json:"id"`
	UserID           uint      `json:"user_id"`
	Name             string    `json:"name"`
	Sitekey          string    `json:"sitekey"`
	Secret           string    `json:"secret"`
	BaseDifficulty   int       `json:"base_difficulty"`
	MaxDifficulty    int       `json:"max_difficulty"`
	ChallengeTimeout int       `json:"challenge_timeout"`
	TokenTimeout     int       `json:"token_timeout"`
	Enabled          bool      `json:"enabled"`
	BatchModeEnabled bool      `json:"batch_mode_enabled"`
	BatchDifficulty  int       `json:"batch_difficulty"`
	MinBatchCount    int       `json:"min_batch_count"`
	MaxBatchCount    int       `json:"max_batch_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	Domains []SiteDomain `json:"domains"`
	Owner   *user.User   `json:"owner,omitempty"`
}

// SiteDomain 站点域名（值对象）
type SiteDomain struct {
	ID        uint      `json:"id"`
	SiteID    uint      `json:"site_id"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"created_at"`
}

// NewSite 创建新站点
func NewSite(userID uint, name string) *Site {
	return &Site{
		UserID:           userID,
		Name:             name,
		Sitekey:          generateSitekey(),
		Secret:           generateSecret(),
		BaseDifficulty:   4,
		MaxDifficulty:    8,
		ChallengeTimeout: 120,
		TokenTimeout:     60,
		Enabled:          true,
		BatchModeEnabled: true,
		BatchDifficulty:  4,
		MinBatchCount:    1,
		MaxBatchCount:    10,
	}
}

// ValidateDomain 验证域名是否在白名单中
func (s *Site) ValidateDomain(origin string) bool {
	if len(s.Domains) == 0 {
		return true
	}

	hostname := extractHostname(origin)
	if hostname == "" {
		return false
	}

	for _, domain := range s.Domains {
		// 存储的域名可能包含端口（如 localhost:5173），
		// 从 origin 提取的 hostname 已剥离端口，所以域名也需要同样处理
		pattern := extractHostname(domain.Domain)
		if matchDomain(pattern, hostname) {
			return true
		}
	}

	return false
}

// RegenerateSecret 重新生成 secret
func (s *Site) RegenerateSecret() {
	s.Secret = generateSecret()
}

func generateSitekey() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")[:12]
}

func generateSecret() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func extractHostname(origin string) string {
	origin = strings.TrimPrefix(origin, "http://")
	origin = strings.TrimPrefix(origin, "https://")
	if idx := strings.Index(origin, ":"); idx != -1 {
		origin = origin[:idx]
	}
	if idx := strings.Index(origin, "/"); idx != -1 {
		origin = origin[:idx]
	}
	return strings.ToLower(origin)
}

func matchDomain(pattern, hostname string) bool {
	pattern = strings.ToLower(pattern)
	hostname = strings.ToLower(hostname)
	if pattern == hostname {
		return true
	}
	if strings.HasPrefix(pattern, "*.") {
		suffix := pattern[1:]
		return strings.HasSuffix(hostname, suffix)
	}
	return false
}

// NewSiteDomain 创建新的站点域名
func NewSiteDomain(siteID uint, domain string) SiteDomain {
	return SiteDomain{
		SiteID: siteID,
		Domain: strings.TrimSpace(domain),
	}
}
