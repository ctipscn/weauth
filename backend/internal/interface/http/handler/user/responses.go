package user

import "time"

// UserResponse 用户响应
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}

// ProfileResponse 用户资料响应
type ProfileResponse struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Avatar      string    `json:"avatar"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}