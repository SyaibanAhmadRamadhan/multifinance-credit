package limits

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (r *repository) Get(ctx context.Context, input GetInput) (output GetOutput, err error) {
	query := r.sq.Select("id", "consumer_id", "tenor", "amount", "remaining_amount").
		From("limits")

	rdbms := r.sqlx
	if input.Tx != nil {
		rdbms = input.Tx
		query = query.Suffix(fmt.Sprintf("FOR %s", input.Locking))
	}

	if input.ID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ID.Int64})
	}
	if input.ConsumerID.Valid {
		query = query.Where(squirrel.Eq{"consumer_id": input.ConsumerID.Int64})
	}
	if input.Tenor.Valid {
		query = query.Where(squirrel.Eq{"tenor": input.Tenor.Int32})
	}

	err = rdbms.QueryRowSq(ctx, query, db.QueryRowScanTypeStruct, &output)
	if err != nil {
		return output, tracer.Error(err)
	}
	return
}
