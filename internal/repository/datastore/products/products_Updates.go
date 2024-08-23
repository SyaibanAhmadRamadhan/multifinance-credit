package products

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (r *repository) Updates(ctx context.Context, input UpdatesInput) (err error) {
	if input.Items == nil || len(input.Items) == 0 {
		return
	}

	rdbms := r.sqlx
	if input.Transaction != nil {
		rdbms = input.Transaction
	}

	caseClause := squirrel.Case("id")

	for _, item := range input.Items {
		caseClause = caseClause.When(fmt.Sprintf("%d", item.ID), fmt.Sprintf("%d", item.Qty))
	}

	rawQuery, args, err := r.sq.Update("products").Set("qty", caseClause).
		Where(squirrel.Eq{"id": getItemIDs(input.Items)}).ToSql()
	if err != nil {
		return tracer.Error(err)
	}

	_, err = rdbms.ExecContext(ctx, rawQuery, args...)
	if err != nil {
		return tracer.Error(err)
	}
	return
}

func getItemIDs(items []UpdatesInputItem) []int64 {
	ids := make([]int64, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	return ids
}
