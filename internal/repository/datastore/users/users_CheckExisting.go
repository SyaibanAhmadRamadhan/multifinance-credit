package users

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (r *repository) CheckExisting(ctx context.Context, input CheckExistingInput) (output CheckExistingOutput, err error) {
	query := r.sq.Select("1").Prefix("SELECT EXISTS(").
		From("users").Suffix(")")
	if input.ByID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ByID.Int64})
	}
	if input.ByEmail.Valid {
		query = query.Where(squirrel.Eq{"email": input.ByEmail.String})
	}

	var existing bool
	err = r.sqlx.QueryRowSq(ctx, query, db.QueryRowScanTypeDefault, &existing)
	if err != nil {
		return output, tracer.Error(err)
	}

	output = CheckExistingOutput{
		Existing: existing,
	}
	return
}
