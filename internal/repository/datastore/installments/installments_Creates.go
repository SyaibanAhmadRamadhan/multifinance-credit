package installments

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

	query := r.sq.Insert("installments").Columns(
		"limit_id", "contract_number", "amount", "due_date", "status",
	)

	for _, item := range input.Items {
		query = query.Values(input.LimitID, input.ContractNumber, item.Amount, item.DueDate, item.Status)
	}

	_, err = rdbms.ExecSq(ctx, query)
	if err != nil {
		return tracer.Error(err)
	}
	return
}
