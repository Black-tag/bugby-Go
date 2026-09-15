package users

import (
	"time"
)

type CreateUserRequest struct {
	Name    string `json:"name"`
	EmailID string `json:"email_id"`
	Password string
	RoleID string `json:"role_id"`
}

type CreateUserResponse struct {
	Name      string
	EmailID   string
	CreatedAt time.Time
	Role      string
}
