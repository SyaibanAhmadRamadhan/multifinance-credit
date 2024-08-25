package products

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (r *repository) Get(ctx context.Context, input GetInput) (output GetOutput, err error) {
	query := r.sq.Select("id", "merchant_id", "name", "image", "qty", "price").
		From("products")

	if input.ID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ID.Int64})
	}

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return output, tracer.Error(err)
	}

	row := r.sqlx.QueryRowxContext(ctx, rawQuery, args...)

	err = row.StructScan(&output)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Join(err, datastore.ErrRecordNotFound)
		}
		return output, tracer.Error(err)
	}
	return
}
