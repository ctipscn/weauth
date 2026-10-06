package user

// UserQuery 用户查询过滤条件
type UserQuery struct {
	ID       uint
	Username string
	Email    string
	Role     Role
	Enabled  *bool
}