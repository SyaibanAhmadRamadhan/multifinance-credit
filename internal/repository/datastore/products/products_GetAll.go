package products

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/rs/zerolog/log"
)

func (r *repository) GetAll(ctx context.Context, input GetAllInput) (output GetAllOutput, err error) {
	query := r.sq.Select("id", "merchant_id", "name", "image", "qty", "price").
		From("products")
	queryCount := r.sq.Select("COUNT(*)").From("products")

	if input.MerchantID.Valid {
		query = query.Where(squirrel.Eq{"merchant_id": input.MerchantID.Int64})
		queryCount = queryCount.Where(squirrel.Eq{"merchant_id": input.MerchantID.Int64})
	}
	if input.IDs != nil {
		query = query.Where(squirrel.Eq{"id": input.IDs})
		queryCount = queryCount.Where(squirrel.Eq{"id": input.IDs})
	}

	rawQueryCount, args, err := queryCount.ToSql()
	if err != nil {
		return output, tracer.Error(err)
	}
	totalData := int64(0)

	row, err := r.sqlx.QueryRowxContext(ctx, rawQueryCount, args...)
	if err != nil {
		return output, tracer.Error(err)
	}

	err = row.Scan(&totalData)
	if err != nil {
		return output, tracer.Error(err)
	}

	offset := pagination.GetOffsetValue(input.Pagination.Page, input.Pagination.PageSize)
	query = query.Limit(uint64(input.Pagination.PageSize))
	query = query.Offset(uint64(offset))

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return output, tracer.Error(err)
	}

	rows, err := r.sqlx.QueryxContext(ctx, rawQuery, args...)
	if err != nil {
		return output, tracer.Error(err)
	}
	defer func() {
		if errRowsClose := rows.Close(); errRowsClose != nil {
			log.Err(errRowsClose).Msg("failed closed row")
		}
	}()

	output = GetAllOutput{
		Pagination: pagination.CreatePaginationOutput(input.Pagination, totalData),
		Items:      make([]GetAllOutputItem, 0),
	}

	for rows.Next() {
		item := GetAllOutputItem{}
		err = rows.StructScan(&item)
		if err != nil {
			return output, tracer.Error(err)
		}

		output.Items = append(output.Items, item)
	}
	return
}
