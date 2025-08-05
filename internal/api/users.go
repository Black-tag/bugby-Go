package api

import (
	"database/sql"
	"encoding/json"
	

	"log"
	"net/http"
	"time"

	"github.com/blacktag/bugby-Go/internal/database"
	"github.com/blacktag/bugby-Go/internal/utils"
	"github.com/google/uuid"
)
type CreateUserRequest struct {
		Email string `json:"email" example:"user@example.com"`
		Password string `json:"password" example:"mysecret"`
		
		
}

type CreateUserResponse struct {
        ID        uuid.UUID `json:"id"`
        Email     string    `json:"email"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
		
    }

type LoginUserRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"mysecret"`
		
}
type LoginResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}
type UpdateRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"mysecret"`
}

type UpdateResponse struct {
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


// @Summary Creates a new  user
// @Description Creates user with Email and Password
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User creation data"
// @Success 201 {object} CreateUserResponse
// @Failure 400 {object} utils.ErrorResponse "Bad Request - Invalid input"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /users [post]
// @Security BearerAuth
func (cfg *APIConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")
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
		log.Printf("cannot create user: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot create user")
		return
	 }
	 utils.RespondWithJSON(w, http.StatusCreated, CreateUserResponse{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	 })

}
// @Summary Login an existing  user
// @Description Existing users can login using email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body LoginUserRequest true "user login data"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} utils.ErrorResponse "Bad Request - Invalid input"
// @Failure 401 {object} utils.ErrorResponse "Invalid credentials" example:{"code":401,"message":"invalid email or password"}
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /login [post]
// @Security BearerAuth
func (cfg *APIConfig) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	
	var req LoginUserRequest
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
	
	token, err := utils.MakeJWT(user.ID, cfg.SECRET,time.Hour)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot generate token")
		return
		

	}
	refreshToken, err := utils.MakeRefreshToken()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot refresh token")
	}
	refreshExpiresAt := time.Now().Add(60 * 24 * time.Hour)

	err = cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ExpiresAt: refreshExpiresAt,
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot store refreshed token")
	}
	utils.RespondWithJSON(w, http.StatusOK, LoginResponse{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Token: token,
		RefreshToken: refreshToken,

	})
}
// @Summary Refresh jwtoken of an existing user
// @Description Existing users can refresh jwt token for future use
// @Tags refreshTokens
// @Accept json
// @Produce json
// @Success 201
// @Failure 400 {object} utils.ErrorResponse "Bad Request - Invalid input"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Missing/invalid credentials"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /refresh [post]
// @Security BearerAuth
func (cfg *APIConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := utils.GetBearerToken(r.Header)
	
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}
	refreshToken, err := cfg.DB.GetUserFromRefreshToken(r.Context(),tokenString)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}
	
	if refreshToken.RevokedAt.Valid || time.Now().After(refreshToken.ExpiresAt){ 
		utils.RespondWithError(w, http.StatusUnauthorized, "refresh token expired or revoked")
		return
	}
	newAccessToken, err := utils.MakeJWT(refreshToken.UserID, cfg.SECRET, time.Hour)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError,"failed to create access token")
		return 
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"token": newAccessToken,
	})	
	
	
}
// @Summary Revoke user token
// @Description Existing users can revoke token using email and password
// @Tags users
// @Accept json
// @Produce json
// @Success 204 
// @Failure 400 {object} utils.ErrorResponse "Bad Request - Invalid input"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Missing/invalid credentials"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /revoke [post]
// @Security BearerAuth
func (cfg *APIConfig) RevokeTokenHandler (w http.ResponseWriter, r *http.Request) {
	tokenString, ok := r.Context().Value("refreshTokenString").(string)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid token string in context")
		return
	}
	err := cfg.DB.RevokeRefreshToken(r.Context(), tokenString)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to revoke token")
        return
    }
	w.WriteHeader(http.StatusNoContent)
}


// @Summary Update an existing  user
// @Description Existing users can update their info using email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body UpdateRequest true "User updation data"
// @Success 201 {object} UpdateResponse
// @Failure 400 {object} utils.ErrorResponse "Bad Request - Invalid input"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Missing/invalid credentials"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /users [put]
// @Security BearerAuth
func (cfg *APIConfig) UpdateCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	
	userIDVal := r.Context().Value("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid or missing user ID")
		return
	}

	
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	hashedpassword, err :=utils.HashPassword(req.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot hash the new password")
		return
	}
	err = cfg.DB.UpdateUserCredentials(r.Context(),database.UpdateUserCredentialsParams{
		Email: req.Email,
		HashedPassword: hashedpassword,
		ID: userID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to update user credentials")
		return
	}
	user, err := cfg.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch the updated user")
		return 
	}
	
	utils.RespondWithJSON(w, http.StatusOK, UpdateResponse{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		
		
	})
}






