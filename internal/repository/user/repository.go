package user

import (
	"github.com/s0vunia/auth_microservice/internal/repository"
	"github.com/s0vunia/platform_common/pkg/db"
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
