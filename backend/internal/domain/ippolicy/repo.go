package ippolicy

// Repository IP 策略仓库接口
type Repository interface {
	FindBySiteID(siteID uint) (*IpPolicy, error)
	FindGlobal() (*IpPolicy, error)
	CreateOrUpdate(siteID uint, policy *IpPolicy) (*IpPolicy, error)
}