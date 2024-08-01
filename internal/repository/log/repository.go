package log

import (
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/repository"
	"github.com/s0vunia/platform_common/pkg/db"
)

var (
	tableName = "logs"

	_             = "id"
	messageColumn = "message"
	_             = "created_at"
)

type repo struct {
	db db.Client
}

// NewRepository creates a new log repository.
func NewRepository(db db.Client) repository.LogRepository {
	return &repo{db: db}
}
