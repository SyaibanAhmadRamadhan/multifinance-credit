package consumers

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"time"
)

func (r *repository) Create(ctx context.Context, input CreateInput) (output CreateOutput, err error) {
	rdbms := r.sqlx
	if input.Transaction != nil {
		rdbms = input.Transaction
	}

	query := r.sq.Insert("consumers").
		Columns("user_id", "nik", "full_name", "legal_name", "place_of_birth", "date_of_birth", "salary",
			"photo_ktp", "photo_selfie", "created_at", "updated_at").
		Values(
			input.UserID, input.Nik, input.FullName, input.LegalName,
			input.PlaceOfBirth, input.DateOfBirth, input.Salary, input.PhotoKTP, input.PhotoSelfie,
			time.Now().UTC(), time.Now().UTC(),
		)

	res, err := rdbms.ExecSq(ctx, query)
	if err != nil {
		return output, tracer.Error(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return output, tracer.Error(err)
	}

	output = CreateOutput{
		ID: id,
	}
	return
}
