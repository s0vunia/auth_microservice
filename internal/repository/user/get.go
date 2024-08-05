package user

import (
	"context"

	"github.com/s0vunia/platform_common/pkg/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/s0vunia/auth_microservice/internal/model"
	"github.com/s0vunia/auth_microservice/internal/repository/user/converter"
	modelRepo "github.com/s0vunia/auth_microservice/internal/repository/user/model"
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
