package admin

// DeleteUserRequest 删除用户请求体
type DeleteUserRequest struct {
	CascadeSites bool `json:"cascade_sites"` // 是否级联删除该用户下的所有站点
}