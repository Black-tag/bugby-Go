package utils

import (
	
	
	"fmt"
	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (string, error) {
	hashedPassword , err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	fmt.Println("hashedpassword", string(hashedPassword))
	return string(hashedPassword) ,nil
}



func CheckPasswordAndHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}