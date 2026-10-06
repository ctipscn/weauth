package verifylog

// Repository 验证日志仓储接口
type Repository interface {
	Create(log *VerificationLog) error
	CountBySiteID(siteID uint) (int64, error)
	CountSuccessBySiteID(siteID uint) (int64, error)
}