package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/blacktag/bugby-Go/internal/database"
	"github.com/blacktag/bugby-Go/internal/utils"
	"github.com/google/uuid"
)
type CreateUserRequest struct {
		Email string `json:"email"`
		Password string `json:"password"`
		
}

type createUserResponse struct {
        ID        uuid.UUID `json:"id"`
        Email     string    `json:"email"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
		
    }



func (cfg *APIConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	

	var req CreateUserRequest
	 err := json.NewDecoder(r.Body).Decode(&req)
	 if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	 }

	 if req.Email == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "email field required")
		return
	 }
	 hashed_password, err := utils.HashPassword(req.Password)
	 if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot hashPassword")
		return
	 }

	 user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email: req.Email,
		HashedPassword: hashed_password,
	 })
	 if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot create user")
		return
	 }
	 utils.RespondWithJSON(w, http.StatusCreated, createUserResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	 })

}

func (cfg *APIConfig) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	type loginuserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type loginResponse struct {
		ID           uuid.UUID `json:"id"`
		Email        string    `json:"email"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}
	var req loginuserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request budy")
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "user does not exists")
		return
	}

	

	err = utils.CheckPasswordAndHash(req.Password, user.HashedPassword )
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, loginResponse{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}