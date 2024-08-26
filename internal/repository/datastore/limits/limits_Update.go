package limits

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (r *repository) Update(ctx context.Context, input UpdateInput) (err error) {
	query := r.sq.Update("limits").
		Set("remaining_amount", input.RemainingAmount).
		Where(squirrel.Eq{"id": input.ID})

	rdbms := r.sqlx
	if input.Transaction != nil {
		rdbms = input.Transaction
	}

	_, err = rdbms.ExecSq(ctx, query)
	if err != nil {
		return tracer.Error(err)
	}
	return
}
