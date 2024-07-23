package user

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/s0vunia/auth_microservices_course_boilerplate/internal/client/db"
)

func (r *repo) IsExists(ctx context.Context, ids []int64) (bool, error) {
	// Convert the slice of int64 to a slice of interfaces{}, which is required by Squirrel
	idArgs := make([]interface{}, len(ids))
	for i, id := range ids {
		idArgs[i] = id
	}

	builder := sq.Select(fmt.Sprintf("COUNT(DISTINCT %s)", idColumn)).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: idArgs})

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

	return count == len(ids), nil
}
