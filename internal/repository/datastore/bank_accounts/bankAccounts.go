package bank_accounts

import (
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
)

type repository struct {
	sqlx *db.SqlxWrapper
	sq   squirrel.StatementBuilderType
}

func NewRepository(sqlx *db.SqlxWrapper) *repository {
	return &repository{
		sqlx: sqlx,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}
