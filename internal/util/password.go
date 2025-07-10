package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(base string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(base), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("there was an error hashing the given password. Cause: %s", err.Error())
	}

	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
