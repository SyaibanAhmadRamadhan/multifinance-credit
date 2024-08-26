package installments

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/jmoiron/sqlx"
)

func (r *repository) GetAll(ctx context.Context, input GetAllInput) (output GetAllOutput, err error) {
	query := r.sq.Select("id", "limit_id", "contract_number", "amount", "due_date", "payment_date", "status").
		From("installments").
		Where(squirrel.Eq{"contract_number": input.ContractNumber}).
		OrderBy("due_date ASC").
		OrderBy("id DESC")

	queryCount := r.sq.Select("COUNT(*)").From("installments").
		Where(squirrel.Eq{"contract_number": input.ContractNumber})

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
