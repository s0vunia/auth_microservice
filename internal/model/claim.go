package model

import "github.com/dgrijalva/jwt-go"

const (
	ExamplePath = "/user_v1.UserV1/Get"
)

type UserClaims struct {
	jwt.StandardClaims
	ID   int64 `json:"id"`
	Role Role  `json:"role"`
}
