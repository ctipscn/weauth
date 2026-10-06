package iprecord

// Repository IP 记录仓库接口
type Repository interface {
	FindBySiteIDAndIP(siteID uint, clientIP string) (*IpRecord, error)
	Create(record *IpRecord) error
	Save(record *IpRecord) error
	ListBySiteID(siteID uint) ([]*IpRecord, error)
	DeleteBySiteIDAndIP(siteID uint, clientIP string) error
}