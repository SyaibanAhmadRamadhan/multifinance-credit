package limits

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

func DefaultLimitData() []CreatesInputItem {
	return []CreatesInputItem{
		{
			Tenor:  3,
			Amount: 300000,
		},
		{
			Tenor:  6,
			Amount: 600000,
		},
		{
			Tenor:  12,
			Amount: 1800000,
		},
	}
}
