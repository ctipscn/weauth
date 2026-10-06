// Package security 提供密码加密和验证相关的安全服务。
package security


// Encryptor 密码加密接口，提供密码哈希化和比对功能。
type Encryptor interface {
	// HashPassword 对明文密码进行 bcrypt 哈希化，返回哈希字符串。
	HashPassword(password string) (string, error)

	// VerifyPassword 比对明文密码与哈希值是否匹配。
	VerifyPassword(password, hash string) bool

	// GeneratePassword 生成随机强密码，返回密码和可能的错误。
	GeneratePassword() (string, error)
}

