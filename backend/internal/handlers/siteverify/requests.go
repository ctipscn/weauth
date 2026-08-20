package siteverify

// SiteVerifyRequest 站点验证请求
type SiteVerifyRequest struct {
	Secret string `json:"secret"`
	Token  string `json:"token"`
}