package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/blacktag/bugby-Go/internal/database"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (service *UserService) createUser(ctx context.Context, args CreateUserRequest) (database.User, error) {

	hashedPassword := ""
	roleID, err := uuid.Parse(args.RoleID)
	if err != nil {
		return database.User{}, fmt.Errorf("invalid role_id: %w", err)
	}

	fmt.Println("ROLE ID FROM REQUEST:", args.RoleID)
	fmt.Println("PARSED ROLE ID:", roleID)
	params := database.CreateUserParams{
		Email:        args.EmailID,
		Username:     args.UserName,
		PasswordHash: hashedPassword,
		RoleID:       roleID,
	}

	user, err := service.repo.CreateUser(ctx, params)
	if err != nil {
		return database.User{}, err
	}
	return user, err

}
