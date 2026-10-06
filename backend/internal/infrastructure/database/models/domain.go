package models

import "time"

// Domain 站点域名白名单
type Domain struct {
	ID        uint   `gorm:"primaryKey"`
	SiteID    uint   `gorm:"not null;index"`
	Domain    string `gorm:"size:255;not null"`
	CreatedAt time.Time
}

func (Domain) TableName() string {
	return "domains"
}
