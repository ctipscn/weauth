package models

import (
	"time"
)

// IpRecord IP 访问记录表
type IpRecord struct {
	ID                uint      `gorm:"primaryKey"`
	SiteID            uint      `gorm:"not null;uniqueIndex:idx_site_ip"`
	ClientIP          string    `gorm:"size:45;not null;uniqueIndex:idx_site_ip"`
	RequestCount      int       `gorm:"default:0"`
	FailureCount      int       `gorm:"default:0"`
	CurrentDifficulty int       `gorm:"default:4"`
	WindowStart       time.Time `gorm:"not null"`
	LastRequestAt     time.Time
	CreatedAt         time.Time

}

func (IpRecord) TableName() string {
	return "ip_records"
}