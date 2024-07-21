package user

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/client/db"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/model"
)

func (r *repo) Update(ctx context.Context, id int64, userUpdate *model.UserUpdate) error {
	builderUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(nameColumn, userUpdate.Name).
		Set(emailColumn, userUpdate.Email).
		Set(roleColumn, userUpdate.Role).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: id})
	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}
	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}
