package users

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

	query, args, err := r.sq.Insert("users").
		Columns("email", "password", "created_at").
		Values(input.Email, input.Password, time.Now().UTC()).ToSql()
	if err != nil {
		return output, tracer.Error(err)
	}

	res, err := rdbms.Exec(ctx, query, args...)
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
