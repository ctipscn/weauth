package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Site 站点表
type Site struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uint           `gorm:"index;default:0" json:"user_id"` // 站点所有者
	Name             string         `gorm:"size:100;not null" json:"name"`
	Sitekey          string         `gorm:"size:32;uniqueIndex;not null" json:"sitekey"`
	Secret           string         `gorm:"size:64;not null" json:"-"`
	BaseDifficulty   int            `gorm:"default:4;not null" json:"base_difficulty"`
	MaxDifficulty    int            `gorm:"default:8;not null" json:"max_difficulty"`
	ChallengeTimeout int            `gorm:"default:120" json:"challenge_timeout"`
	TokenTimeout     int            `gorm:"default:300" json:"token_timeout"`
	Enabled          bool           `gorm:"default:true" json:"enabled"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// 批量模式配置
	BatchModeEnabled bool `gorm:"default:false" json:"batch_mode_enabled"` // 是否启用批量模式
	BatchDifficulty  int  `gorm:"default:4" json:"batch_difficulty"`       // 批量模式固定难度
	MinBatchCount    int  `gorm:"default:1" json:"min_batch_count"`        // 最小批次数
	MaxBatchCount    int  `gorm:"default:10" json:"max_batch_count"`       // 最大批次数

	// 关联
	Domains  []SiteDomain `gorm:"foreignKey:SiteID" json:"domains,omitempty"`
	IpPolicy *IpPolicy    `gorm:"foreignKey:SiteID" json:"ip_policy,omitempty"`
	Owner    *User        `gorm:"foreignKey:UserID" json:"owner,omitempty"`
}

// SiteDomain 站点域名白名单
type SiteDomain struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SiteID    uint      `gorm:"not null;index" json:"site_id"`
	Domain    string    `gorm:"size:255;not null" json:"domain"`
	CreatedAt time.Time `json:"created_at"`

	// 唯一约束
	Site Site `gorm:"foreignKey:SiteID" json:"-"`
}

// TableName 设置表名
func (Site) TableName() string {
	return "sites"
}

func (SiteDomain) TableName() string {
	return "site_domains"
}

// BeforeCreate 创建前生成 sitekey 和 secret
func (s *Site) BeforeCreate(tx *gorm.DB) error {
	if s.Sitekey == "" {
		s.Sitekey = generateSitekey()
	}
	if s.Secret == "" {
		s.Secret = generateSecret()
	}
	return nil
}

// generateSitekey 生成 sitekey（短而唯一）
func generateSitekey() string {
	// 使用UUID前12个字符（去掉连字符），足够避免冲突
	return strings.ReplaceAll(uuid.New().String(), "-", "")[:12]
}

// generateSecret 生成 secret（使用单个UUID去掉连字符）
func generateSecret() string {
	// UUID去掉连字符后是32个字符，已经足够安全
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// RegenerateSecret 重新生成 secret
func (s *Site) RegenerateSecret() {
	s.Secret = generateSecret()
}
