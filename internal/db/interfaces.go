package db

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type SqlxTransaction interface {
	DoTransaction(ctx context.Context, opt *sql.TxOptions, fn func(tx Rdbms) error) error
}

type Rdbms interface {
	QueryxContext(ctx context.Context, query string, arg ...interface{}) (*sqlx.Rows, error)
	ExecContext(ctx context.Context, query string, arg ...interface{}) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}
