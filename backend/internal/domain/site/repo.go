package site

// Repository 站点仓库接口
type Repository interface {
	FindByID(id uint) (*Site, error)
	FindByIDWithUser(id uint, userID uint) (*Site, error)
	FindAll() ([]*Site, error)
	FindBySitekey(sitekey string) (*Site, error)
	FindBySecret(secret string) (*Site, error)
	ListByUser(userID uint) ([]*Site, error)
	Create(site *Site) error
	Transaction(fn func(repo Repository) error) error
	// 域名管理
	AddDomain(siteID uint, domain string) (*SiteDomain, error)
	DeleteDomain(siteID, domainID uint) error
	// 全量更新（用于事务内）
	Save(site *Site) error
	Delete(id uint) error

	// Unscoped 方法：管理员绕过 GORM 的软删除过滤，查看/操作已软删除的记录
	FindAllUnscoped() ([]*Site, error)
	FindByIDUnscoped(id uint) (*Site, error)
}