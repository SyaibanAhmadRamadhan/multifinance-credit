package consumers

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (r *repository) CheckExisting(ctx context.Context, input CheckExistingInput) (output CheckExistingOutput, err error) {
	query := r.sq.Select("1").Prefix("SELECT EXISTS(").
		From("consumers")
	if input.ByNIK.Valid {
		query = query.Where(squirrel.Eq{"nik": input.ByNIK.String})
	}
	if input.ByID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ByID.Int64})
	}

	sql, args, err := query.Suffix(")").ToSql()
	if err != nil {
		return output, tracer.Error(err)
	}

	row, err := r.sqlx.QueryRowxContext(ctx, sql, args...)
	if err != nil {
		return output, tracer.Error(err)
	}

	var existing bool
	err = row.Scan(&existing)
	if err != nil {
		return output, tracer.Error(err)
	}

	output = CheckExistingOutput{
		Existing: existing,
	}
	return

}
