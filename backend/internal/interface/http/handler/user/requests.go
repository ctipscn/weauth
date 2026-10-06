package user

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}