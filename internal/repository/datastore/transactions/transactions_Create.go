package transactions

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"time"
)

func (r *repository) Create(ctx context.Context, input CreateInput) (output CreateOutput, err error) {
	rdbms := r.sqlx
	if input.Transaction != nil {
		rdbms = input.Transaction
	}

	query := r.sq.Insert("transactions").Columns(
		"limit_id", "consumer_id", "contract_number", "amount", "transaction_date", "status", "created_at", "updated_at",
	).Values(
		input.LimitID, input.ConsumerID, input.ContractNumber, input.Amount, input.TransactionDate, input.Status,
		time.Now().UTC(), time.Now().UTC(),
	)

	res, err := rdbms.ExecSq(ctx, query)
	if err != nil {
		return output, tracer.Error(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return output, tracer.Error(err)
	}

	output = CreateOutput{
		ID: id,
	}
	return
}
