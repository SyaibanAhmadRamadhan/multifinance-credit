package products

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

	query := r.sq.Insert("products").Columns(
		"merchant_id", "image", "name", "qty", "price",
	)

	for _, item := range input.Items {
		query = query.Values(input.MerchantID, item.Image, item.Name, item.Qty, item.Price)
	}

	_, err = rdbms.ExecSq(ctx, query)
	if err != nil {
		return tracer.Error(err)
	}
	return
}
