package model

import (
	"database/sql"
	"time"
)

// User represents a user entity with ID, Info, CreatedAt, and UpdatedAt fields.
type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// UserInfo represents a user info entity with Name, Email, and Role fields.
type UserInfo struct {
	Name  string
	Email string
	Role  Role
}

// UserCreate represents a user create entity with Name, Email, Role, and Password fields.
type UserCreate struct {
	Name     string
	Email    string
	Role     Role
	Password string
}

// UserUpdate represents a user update entity with Name, Email, and Role fields.
type UserUpdate struct {
	Name  *string
	Email *string
	Role  Role
}

// UserLogin represents a user login entity with ID and Password fields.
type UserLogin struct {
	ID       int64
	Password string
}

// Role represents a user role
type Role int32

const (
	// RoleUnknown represents an unknown user role
	RoleUnknown Role = 0
	// RoleUser represents a user role
	RoleUser Role = 1
	// RoleAdmin represents an admin role
	RoleAdmin Role = 2
)
