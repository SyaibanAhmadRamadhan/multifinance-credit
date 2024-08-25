package bank_accounts

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

	query := r.sq.Insert("bank_accounts").Columns(
		"consumer_id", "name", "account_number", "account_holder_name",
	)

	for _, item := range input.Items {
		query = query.Values(item.ConsumerID, item.Name, item.AccountNumber, item.AccountHolderName)
	}

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return tracer.Error(err)
	}

	_, err = rdbms.Exec(ctx, rawQuery, args...)
	if err != nil {
		return tracer.Error(err)
	}
	return
}
