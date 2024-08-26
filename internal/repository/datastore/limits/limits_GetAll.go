package limits

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/jmoiron/sqlx"
)

func (r *repository) GetAll(ctx context.Context, input GetAllInput) (output GetAllOutput, err error) {
	query := r.sq.Select("id", "consumer_id", "tenor", "amount").
		From("limits").OrderBy("id DESC")
	queryCount := r.sq.Select("COUNT(*)").From("limits")
	if input.ConsumerID.Valid {
		query = query.Where(squirrel.Eq{"consumer_id": input.ConsumerID.Int64})
		queryCount = queryCount.Where(squirrel.Eq{"consumer_id": input.ConsumerID.Int64})
	}

	output = GetAllOutput{
		Items: make([]GetAllOutputItem, 0),
	}

	output.Pagination, err = r.sqlx.QuerySqPagination(ctx, queryCount, query, input.Pagination, func(rows *sqlx.Rows) (err error) {
		for rows.Next() {
			item := GetAllOutputItem{}
			err = rows.StructScan(&item)
			if err != nil {
				return tracer.Error(err)
			}

			output.Items = append(output.Items, item)
		}
		return
	})
	if err != nil {
		return output, tracer.Error(err)
	}

	return
}
