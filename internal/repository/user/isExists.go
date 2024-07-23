package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/client/db"
)

func (r *repo) IsExists(ctx context.Context, id int64) (bool, error) {
	builder := sq.Select("COUNT(*)").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return false, err
	}
	q := db.Query{
		Name:     "user_repository.IsExists",
		QueryRaw: query,
	}
	var count int
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
