package model

type User struct {
	ID          int64  `redis:"id"`
	Name        string `redis:"name"`
	Email       string `redis:"email"`
	Role        Role   `redis:"role"`
	CreatedAtNs int64  `redis:"created_at"`
	UpdatedAtNs *int64 `redis:"updated_at"`
}

type Role int32

const (
	RoleUnknown Role = 0
	RoleUser    Role = 1
	RoleAdmin   Role = 2
)
