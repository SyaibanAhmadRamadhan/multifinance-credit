package limits

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (r *repository) Creates(ctx context.Context, input CreatesInput) (err error) {
	if input.Items == nil || len(input.Items) == 0 {
		return
	}

	rdbms := r.sqlx
	if input.Transaction != nil {
		rdbms = input.Transaction
	}

	query := r.sq.Insert("limits").Columns(
		"consumer_id", "tenor", "amount", "remaining_amount",
	)

	for _, item := range input.Items {
		query = query.Values(input.ConsumerID, item.Tenor, item.Amount, item.Amount)
	}

	_, err = rdbms.ExecSq(ctx, query)
	if err != nil {
		return tracer.Error(err)
	}
	return
}
