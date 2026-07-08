package models

import (
	"time"
)

// IpRecord IP 访问记录表
type IpRecord struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	SiteID            uint      `gorm:"not null;uniqueIndex:idx_site_ip" json:"site_id"`
	ClientIP          string    `gorm:"size:45;not null;uniqueIndex:idx_site_ip" json:"client_ip"`
	RequestCount      int       `gorm:"default:0" json:"request_count"`
	FailureCount      int       `gorm:"default:0" json:"failure_count"`
	CurrentDifficulty int       `gorm:"default:4" json:"current_difficulty"`
	WindowStart       time.Time `gorm:"not null" json:"window_start"`
	LastRequestAt     time.Time `json:"last_request_at"`
	CreatedAt         time.Time `json:"created_at"`

	// 关联
	Site Site `gorm:"foreignKey:SiteID" json:"-"`
}

func (IpRecord) TableName() string {
	return "ip_records"
}

// IsWindowExpired 检查时间窗口是否过期
func (r *IpRecord) IsWindowExpired(windowSeconds int) bool {
	return time.Now().After(r.WindowStart.Add(time.Duration(windowSeconds) * time.Second))
}

// ResetWindow 重置统计窗口
func (r *IpRecord) ResetWindow() {
	r.WindowStart = time.Now()
	r.RequestCount = 0
	r.FailureCount = 0
}

// IncrementRequest 增加请求计数
func (r *IpRecord) IncrementRequest() {
	r.RequestCount++
	r.LastRequestAt = time.Now()
}

// IncrementFailure 增加失败计数
func (r *IpRecord) IncrementFailure() {
	r.FailureCount++
	r.LastRequestAt = time.Now()
}
