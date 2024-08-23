package limits

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (r *repository) Update(ctx context.Context, input UpdateInput) (err error) {
	rawQuery, args, err := r.sq.Update("limits").
		Set("remaining_amount", input.RemainingAmount).
		Where(squirrel.Eq{"id": input.ID}).
		ToSql()
	if err != nil {
		return tracer.Error(err)
	}

	rdbms := r.sqlx
	if input.Transaction != nil {
		rdbms = input.Transaction
	}

	_, err = rdbms.ExecContext(ctx, rawQuery, args...)
	if err != nil {
		return tracer.Error(err)
	}
	return
}
