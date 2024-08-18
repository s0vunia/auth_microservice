package model

import "github.com/dgrijalva/jwt-go"

// UserClaims represents a user claims
type UserClaims struct {
	jwt.StandardClaims
	ID   int64 `json:"id"`
	Role Role  `json:"role"`
}
