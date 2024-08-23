package consumers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/rs/zerolog/log"
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

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return output, tracer.Error(err)
	}

	row, stmt, err := r.sqlx.QueryRowxContext(ctx, rawQuery, args...)
	if err != nil {
		return output, tracer.Error(err)
	}
	defer func() {
		if errClose := stmt.Close(); errClose != nil {
			log.Err(errClose).Msg("failed closed stmt")
		}
	}()

	err = row.StructScan(&output)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = datastore.ErrRecordNotFound
		}
		return output, tracer.Error(err)
	}

	return
}
