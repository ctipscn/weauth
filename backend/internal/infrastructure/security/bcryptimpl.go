package security

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const passwordChars = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789!@#$%&*"

type bcryptEncryptor struct{}

// NewEncryptor 创建基于 bcrypt 的密码加密服务实例。
func NewEncryptor() Encryptor {
	return &bcryptEncryptor{}
}

func (v *bcryptEncryptor) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (v *bcryptEncryptor) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (v *bcryptEncryptor) GeneratePassword() (string, error) {
	n := 24
	buf := make([]byte, n)
	for i := range(n) {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordChars))))
		if err != nil {
			return "", err
		}
		buf[i] = passwordChars[idx.Int64()]
	}
	return string(buf), nil
}
