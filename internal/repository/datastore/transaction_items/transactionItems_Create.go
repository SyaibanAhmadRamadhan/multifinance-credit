package transaction_items

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (r *repository) Creates(ctx context.Context, input CreatesInput) (err error) {
	rdbms := r.sqlx
	if input.Transaction != nil {
		rdbms = input.Transaction
	}

	query := r.sq.Insert("transaction_items").Columns(
		"transaction_id", "merchant_id", "name", "image", "qty", "unit_price", "amount",
	)

	for _, item := range input.Items {
		query = query.Values(input.TransactionID, item.MerchantID, item.Name, item.Image, item.Qty, item.UnitPrice, item.Amount)
	}

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return tracer.Error(err)
	}

	_, err = rdbms.ExecContext(ctx, rawQuery, args...)
	if err != nil {
		return tracer.Error(err)
	}

	return
}
