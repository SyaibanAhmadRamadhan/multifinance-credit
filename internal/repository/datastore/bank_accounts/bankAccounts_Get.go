package bank_accounts

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (r *repository) Get(ctx context.Context, input GetInput) (output GetOutput, err error) {
	query := r.sq.Select("id", "consumer_id", "name", "account_number", "account_holder_name").
		From("bank_accounts")

	if input.AccountNumber.Valid {
		query = query.Where(squirrel.Eq{"account_number": input.AccountNumber.String})
	}
	if input.ID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ID.Int64})
	}

	err = r.sqlx.QueryRowSq(ctx, query, db.QueryRowScanTypeStruct, &output)
	if err != nil {
		return output, tracer.Error(err)
	}
	return
}
