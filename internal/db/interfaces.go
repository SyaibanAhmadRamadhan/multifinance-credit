package db

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/jmoiron/sqlx"
)

type QueryRowScanType uint8

const (
	QueryRowScanTypeDefault QueryRowScanType = iota + 1
	QueryRowScanTypeStruct
)

type SqlxTransaction interface {
	DoTransaction(ctx context.Context, opt *sql.TxOptions, fn func(tx Rdbms) error) error
}

type callbackRows func(rows *sqlx.Rows) (err error)
type Rdbms interface {
	// reader command

	Query(ctx context.Context, query squirrel.SelectBuilder, callback callbackRows) error
	QueryPagination(ctx context.Context, countQuery, query squirrel.SelectBuilder, pagination pagination.PaginationInput, callback callbackRows) (
		pagination.PaginationOutput, error)
	QueryRow(ctx context.Context, query squirrel.SelectBuilder, scanType QueryRowScanType, dest interface{}) error

	// writer command

	ExecContext(ctx context.Context, query string, arg ...interface{}) (sql.Result, error)
}

type queryExecutor interface {
	QueryxContext(ctx context.Context, query string, arg ...interface{}) (*sqlx.Rows, error)
	ExecContext(ctx context.Context, query string, arg ...interface{}) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}
