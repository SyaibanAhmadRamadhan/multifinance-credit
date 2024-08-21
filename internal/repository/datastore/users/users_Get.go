package users

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (r *repository) Get(ctx context.Context, input GetInput) (output GetOutput, err error) {
	query := r.sq.Select("id", "email", "password").From("users")
	if input.Email.Valid {
		query = query.Where(squirrel.Eq{"email": input.Email.String})
	}
	if input.ID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ID.Int64})
	}

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return output, tracer.Error(err)
	}

	row, err := r.sqlx.QueryRowxContext(ctx, rawQuery, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = datastore.ErrRecordNotFound
		}
		return output, tracer.Error(err)
	}

	err = row.StructScan(&output)
	if err != nil {
		return output, tracer.Error(err)
	}
	return
}
