package users

import (
	"time"
)

type CreateUserRequest struct {
	UserName string `json:"username"`
	EmailID  string `json:"email"`
	Password string `json:"password"`
	RoleID   string `json:"role_id"`
}

type CreateUserResponse struct {
	Name      string
	EmailID   string
	CreatedAt time.Time
	Role      string
}
