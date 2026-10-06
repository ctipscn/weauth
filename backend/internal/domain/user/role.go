package user

// 预计不会出现新的角色，故此处直接在 User 实体中建立内置枚举
type Role int

const (
	Administrator Role = iota
	Staff
)

// String 返回角色的字符串表示（用于 JSON 序列化和数据库交互）
func (r Role) String() string {
	switch r {
	case Administrator:
		return "admin"
	case Staff:
		return "staff"
	default:
		return "staff"
	}
}

// ParseRole 将字符串解析为角色枚举
func ParseRole(s string) (Role, error) {
	switch s {
	case "admin":
		return Administrator, nil
	case "staff":
		return Staff, nil
	default:
		return Staff, nil
	}
}