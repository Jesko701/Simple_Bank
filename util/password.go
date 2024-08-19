package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to generate hash password: %v", err)
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashed_password, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(pw))
}
