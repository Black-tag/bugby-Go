package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/blacktag/bugby-Go/internal/database"
	"github.com/blacktag/bugby-Go/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RefreshTokenFetcher interface {
	GetRefreshToken(ctx context.Context, token string) (database.RefreshToken, error)
}

func Authenticate(secret string, db *database.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "no header")
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 {
				utils.RespondWithError(w, http.StatusUnauthorized, "malfromed troken")
				return
			}
			if parts[0] != "Bearer" {
				utils.RespondWithError(w, http.StatusUnauthorized, "header must contain bearer")
				return
			}
			tokenSring := parts[1]
			claims := &jwt.RegisteredClaims{}
			token, err := jwt.ParseWithClaims(tokenSring, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "error parsing token")
				return
			}
			if !token.Valid {
				utils.RespondWithError(w, http.StatusUnauthorized, "invalid token")
				return
			}
			if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
				utils.RespondWithError(w, http.StatusUnauthorized, "token has expired")
				return
			}
			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "error parsing user id")
				return
			}

			role, err := db.GetRoleByID(r.Context(), userID)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "unable to fetch role")
				return
			}
			ctx := context.WithValue(r.Context(), "userID", userID)
			ctx = context.WithValue(ctx, "tokenString", tokenSring)
			ctx = context.WithValue(ctx, "role", role)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func RevokeTokenAthenticate(db RefreshTokenFetcher) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
				utils.RespondWithError(w, http.StatusUnauthorized, "missing or invalid auth header")
				return
			}
			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.RespondWithError(w, http.StatusUnauthorized, "malformed token")
				return
			}
			tokenString := parts[1]

			refreshToken, err := db.GetRefreshToken(r.Context(), tokenString)

			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "invalid refresh token")
				return
			}
			if refreshToken.RevokedAt.Valid {
				utils.RespondWithError(w, http.StatusUnauthorized, "refresh token revoked")
				return
			}
			if time.Now().After(refreshToken.ExpiresAt) {
				utils.RespondWithError(w, http.StatusUnauthorized, "refresh token expired")
				return
			}
			ctx := context.WithValue(r.Context(), "refreshTokenString", tokenString)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
