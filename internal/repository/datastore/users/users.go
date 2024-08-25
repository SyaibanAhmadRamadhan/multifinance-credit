package users

import (
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
)

type repository struct {
	sqlx db.Rdbms
	sq   squirrel.StatementBuilderType
}

func NewRepository(sqlx db.Rdbms) *repository {
	return &repository{
		sqlx: sqlx,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}
