package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	return string(hashed), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// func comparePasswords(hashedPwd string, plainPwd string) bool {
// 	byteHash := []byte(hashedPwd)
// 	plainPwdBytes := []byte(plainPwd)
// 	err := bcrypt.CompareHashAndPassword(byteHash, plainPwdBytes)
// 	return err == nil
// }

