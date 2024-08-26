package consumers

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (r *repository) Get(ctx context.Context, input GetInput) (output GetOutput, err error) {
	query := r.sq.Select("id", "user_id", "nik", "full_name", "legal_name", "place_of_birth", "date_of_birth",
		"salary", "photo_ktp", "photo_selfie").From("consumers")
	if input.UserID.Valid {
		query = query.Where(squirrel.Eq{"user_id": input.UserID.Int64})
	}
	if input.ID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ID.Int64})
	}

	err = r.sqlx.QueryRowSq(ctx, query, db.QueryRowScanTypeStruct, &output)
	if err != nil {
		return output, tracer.Error(err)
	}

	return
}
