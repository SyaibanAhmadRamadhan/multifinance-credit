package bank_accounts

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
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

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return output, tracer.Error(err)
	}

	row, err := r.sqlx.QueryRowxContext(ctx, rawQuery, args...)
	if err != nil {
		return output, tracer.Error(err)
	}

	err = row.StructScan(&output)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Join(err, datastore.ErrRecordNotFound)
		}
		return output, tracer.Error(err)
	}
	return
}
