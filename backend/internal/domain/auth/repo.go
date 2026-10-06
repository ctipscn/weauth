package auth

type Repository interface {
	FindByUserID(userID uint) (*Auth, error)
	FindByUsername(username string) (*Auth, error)
	FindByOIDCSub(sub string) (*Auth, error)
	Create(credential *Auth) error
	Save(credential *Auth) error
	// DeleteByUserID 软删除指定用户的认证凭证
	DeleteByUserID(userID uint) error
}
