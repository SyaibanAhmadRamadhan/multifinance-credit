package db_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_sqlxTransaction_DoTransaction(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewSqlxTransaction(sqlxDB)

	t.Run("should be return with commit db", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()

		err = sqlxx.DoTransaction(ctx, &sql.TxOptions{}, func(tx db.Rdbms) (err error) {
			return nil
		})
		require.NoError(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should be return with rollback db", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()

		err = sqlxx.DoTransaction(ctx, &sql.TxOptions{}, func(tx db.Rdbms) (err error) {
			return errors.New("error")
		})
		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
