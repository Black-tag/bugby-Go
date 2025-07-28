package api

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"log"
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
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
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
	utils.RespondWithJSON(w, http.StatusOK, loginResponse{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Token: token,
		RefreshToken: refreshToken,

	})
}

func (cfg *APIConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := utils.GetBearerToken(r.Header)
	// fmt.Println("tokenstring:", tokenString)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}
	refreshToken, err := cfg.DB.GetUserFromRefreshToken(r.Context(),tokenString)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}
	// fmt.Println("refreshToken:",refreshToken)
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
func (cfg *APIConfig) RevokeTokenHandler (w http.ResponseWriter, r *http.Request) {
	// tokenString, err := utils.GetBearerToken(r.Header)
	// if err != nil {
	// 	utils.RespondWithError(w, http.StatusUnauthorized, "Missing refresh token")
    //     return
	// }
	// userIDVal := r.Context().Value("userID")
	// userID, ok := userIDVal.(uuid.UUID)
	// if !ok {
	// 	utils.RespondWithError(w, http.StatusUnauthorized, "invalid or missing user ID")
	// 	return
	// }
	// tokenStringVal := r.Context().Value("tokenString")
	// tokenString, ok := tokenStringVal.(string)
	// if !ok {
	// 	utils.RespondWithError(w, http.StatusUnauthorized, "invalid token string in context")
	// 	return
	// }
	tokenString, ok := r.Context().Value("refreshTokenString").(string)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid token string in context")
		return
	}
	


	fmt.Printf("tokenString in revoke handler from context: %v", tokenString)
	err := cfg.DB.RevokeRefreshToken(r.Context(), tokenString)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to revoke token")
        return
    }
	w.WriteHeader(http.StatusNoContent)
}


func (cfg *APIConfig) UpdateCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	
	userIDVal := r.Context().Value("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid or missing user ID")
		return
	}

	type updateReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req updateReq
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
	type UpdateResponse struct {
		ID uuid.UUID `json:"id"`
		Email string `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	utils.RespondWithJSON(w, http.StatusOK, UpdateResponse{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		
		
	})
}






