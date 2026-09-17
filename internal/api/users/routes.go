package users

import (
	"net/http"

	"github.com/blacktag/bugby-Go/internal/database"
	"github.com/blacktag/bugby-Go/internal/validation"
)

func RegisterRoutes(mux *http.ServeMux, db *database.Queries) {
	userRepo := NewRepository(db)
	userService := NewUserService(userRepo)
	userHandler := NewUserHandler(userService, validation.NewValidator())

	mux.HandleFunc("POST /api/v2/users", userHandler.CreateUserHandler)
}
