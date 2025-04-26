package auth

import (
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	fmt.Printf("Checking hash %s against password %s\n", hash, password)
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}