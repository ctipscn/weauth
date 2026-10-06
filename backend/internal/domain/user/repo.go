package user

// Repository 用户仓库接口
type Repository interface {
	CountAll() (int64, error)
	FindAll() ([]*User, error)
	FindByID(id uint) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByEmail(email string) (*User, error)
	Create(user *User) error
	Save(user *User) error
	Delete(id uint) error
	// Deactivate 软删除用户：设置 enabled=false 并填充 deleted_at
	Deactivate(id uint) error
}
