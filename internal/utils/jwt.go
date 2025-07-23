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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.RegisteredClaims{
		Issuer: "bugby",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	})
	signingKey := []byte(tokenSecret)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}
	return signedToken, nil
		

	
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token)(interface{}, error){
		return []byte(tokenSecret), nil

	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing Token")
	}
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return uuid.Nil, fmt.Errorf("token has expired")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing user id")
	}
	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no header")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", errors.New("malformed header ")
	}
	if parts[0] != "Bearer" {
		return "", errors.New("authorization header must contain Bearer")
	}
	return parts[1], nil
}

func MakeRefreshToken () (string, error) {
	data:= make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		fmt.Println("error generating random data:", err)
		return "", err
	}
	encodedData := hex.EncodeToString(data)
	return encodedData, nil
}


func GetAPIKey(headers http.Header)(string, error){
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no header")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", errors.New("malformed header ")
	}
	if parts[0] != "ApiKey" {
		return "", errors.New("authorization header must start with ApiKey")
	}
	return parts[1], nil


}