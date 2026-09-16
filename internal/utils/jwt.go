package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// MakeJWT creates a signed JWT using HS256 and the provided secret.
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "bugby",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})
	signingKey := []byte(tokenSecret)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}
	return signedToken, nil
}

// ValidateJWT parses and validates a token, returning the user UUID if valid.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing token: %w", err)
	}
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return uuid.Nil, fmt.Errorf("token has expired")
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing user id: %w", err)
	}
	return userID, nil
}

// GetBearerToken extracts a bearer token from HTTP headers.
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no header")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", errors.New("malformed header")
	}
	if parts[0] != "Bearer" {
		return "", errors.New("authorization header must contain Bearer")
	}
	return parts[1], nil
}

// GetAPIKey extracts an ApiKey token from Authorization header.
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no header")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", errors.New("malformed header")
	}
	if parts[0] != "ApiKey" {
		return "", errors.New("authorization header must start with ApiKey")
	}
	return parts[1], nil
}

// MakeRefreshToken generates a random 32-byte hex string to use as a refresh token.
func MakeRefreshToken() (string, error) {
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		return "", fmt.Errorf("error generating random data: %w", err)
	}
	return hex.EncodeToString(data), nil
}
