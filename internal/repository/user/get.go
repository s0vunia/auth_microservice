package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/client/db"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/model"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/repository/user/converter"
	modelRepo "github.com/s0vunia/auth_microservices_course_boilerplate/internal/repository/user/model"
)

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {

	builderSelectOne := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})
	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		return nil, err
	}
	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}
	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, err
	}
	return converter.ToUserFromRepo(&user), nil
}
