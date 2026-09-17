package users

import (
	"time"

)

type CreateUserRequest struct {
	UserName string `json:"username" validate:"required"`
	EmailID  string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	RoleID   string `json:"role_id" validate:"required,uuid"`
}

type CreateUserResponse struct {
	Name      string
	EmailID   string
	CreatedAt time.Time
	Role      string
}
