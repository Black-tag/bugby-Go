package users

import (
	"context"
    "github.com/blacktag/bugby-Go/internal/database"
)

type UserRepository interface {
    CreateUser(ctx context.Context, args database.CreateUserParams) (database.User, error)
}

type userRepository struct {
    db *database.Queries
}

func NewRepository(db *database.Queries) UserRepository {
    return &userRepository{
        db: db,
    }
}

func (repo *userRepository) CreateUser(ctx context.Context, args database.CreateUserParams) (database.User, error) {
    // use repo.db here
	return repo.db.CreateUser(ctx, args)
}