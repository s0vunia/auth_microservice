package user

import (
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/client/db"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/repository"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	passHashColumn  = "pass_hash"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

// NewRepository creates a new user repository.
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}
