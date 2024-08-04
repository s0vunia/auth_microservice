package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/client/db"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/model"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/repository/user/converter"
)

func (r *repo) Create(ctx context.Context, userCreate *model.UserCreate) (int64, error) {
	userCreateRepo := converter.ToUserCreateFromService(userCreate)

	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passHashColumn, roleColumn).
		Values(userCreateRepo.Name, userCreateRepo.Email, userCreateRepo.PasswordHash, userCreateRepo.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
