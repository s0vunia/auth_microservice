package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// GetHashPassword generate hash for password
func GetHashPassword(password string) ([]byte, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return make([]byte, 0), err
	}
	return passHash, nil
}

// CompareHashAndPassword compare hash with password
func CompareHashAndPassword(passHash []byte, pass string) error {
	return bcrypt.CompareHashAndPassword(passHash, []byte(pass))
}
