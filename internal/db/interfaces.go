package db

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type SqlxTransaction interface {
	DoTransaction(ctx context.Context, opt *sql.TxOptions, fn func(tx *SqlxWrapper) error) error
}

type queryExecutor interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}
