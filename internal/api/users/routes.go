package users

import (
	"net/http"

	"github.com/blacktag/bugby-Go/internal/database"
)

func RegisterRoutes(mux *http.ServeMux, db *database.Queries) {
	userRepo := NewRepository(db)
	userService := NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	mux.HandleFunc("POST /api/v2/users", userHandler.CreateUserHandler)
}