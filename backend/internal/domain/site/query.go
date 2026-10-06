package site

// SiteQuery 站点查询过滤条件
type SiteQuery struct {
	UserID  uint
	Enabled *bool
	Name    string
}