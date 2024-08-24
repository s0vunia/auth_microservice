package utils

import "golang.org/x/crypto/bcrypt"

// VerifyPassword compares candidate password with hashed password
// returns true if passwords match
func VerifyPassword(hashedPassword string, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}

// HashPassword hashes password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
