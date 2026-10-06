package siteverify

// SiteVerifyRequest 站点验证请求
type SiteVerifyRequest struct {
	Secret string `json:"secret" binding:"required"`
	Token  string `json:"token" binding:"required"`
}