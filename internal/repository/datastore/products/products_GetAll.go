package products

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/jmoiron/sqlx"
)

func (r *repository) GetAll(ctx context.Context, input GetAllInput) (output GetAllOutput, err error) {
	query := r.sq.Select("id", "merchant_id", "name", "image", "qty", "price").
		From("products").OrderBy("id DESC")

	queryCount := r.sq.Select("COUNT(*)").From("products")

	rdbms := r.sqlx
	if input.Transaction != nil {
		rdbms = input.Transaction
		query = query.Suffix(fmt.Sprintf("FOR %s", input.Locking))
		queryCount = queryCount.Suffix(fmt.Sprintf("FOR %s", input.Locking))
	}

	if input.MerchantID.Valid {
		query = query.Where(squirrel.Eq{"merchant_id": input.MerchantID.Int64})
		queryCount = queryCount.Where(squirrel.Eq{"merchant_id": input.MerchantID.Int64})
	}
	if input.IDs != nil {
		query = query.Where(squirrel.Eq{"id": input.IDs})
		queryCount = queryCount.Where(squirrel.Eq{"id": input.IDs})
	}

	output = GetAllOutput{
		Items: make([]GetAllOutputItem, 0),
	}

	output.Pagination, err = rdbms.QuerySqPagination(ctx, queryCount, query, input.Pagination, func(rows *sqlx.Rows) (err error) {
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
