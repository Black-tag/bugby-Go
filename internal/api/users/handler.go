package users

import (
	"encoding/json"
	"net/http"

	"github.com/blacktag/bugby-Go/internal/validation"
// 	"gopkg.in/go-playground/validator.v9"
)

type UserHandler struct {
	service *UserService
	validator *validation.Validator
}

func NewUserHandler(service *UserService, validator *validation.Validator) *UserHandler {
	return &UserHandler{
		service: service,
		validator: validator,
	}
}


func (handler *UserHandler) CreateUserHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err = handler.validator.Validate(req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}

	user, err := handler.service.createUser(r.Context(), req)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}


// func(handler *UserHandler)GetUserWithID(w http.ResponseWriter, r *http.Request){
// 	userID := r.Body()

// }