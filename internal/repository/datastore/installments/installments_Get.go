package installments

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (r *repository) Get(ctx context.Context, input GetInput) (output GetOutput, err error) {
	query := r.sq.Select("id", "limit_id", "contract_number", "amount", "due_date", "payment_date", "status").
		From("installments")

	if input.ID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ID.Int64})
	}

	err = r.sqlx.QueryRow(ctx, query, db.QueryRowScanTypeStruct, &output)
	if err != nil {
		return output, tracer.Error(err)
	}
	return
}
