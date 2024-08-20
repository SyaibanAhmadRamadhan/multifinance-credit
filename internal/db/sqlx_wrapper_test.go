package db_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func Test_sqlxWrapper_Queryx(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewSqlxWrapper(sqlxDB)

	t.Run("should return correct  Queryx result", func(t *testing.T) {
		expectedQuery := `SELECT * FROM "users" WHERE id = ?`

		mock.ExpectPrepare(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = ?`)).
			ExpectQuery().
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		res, err := sqlxx.QueryxContext(ctx, expectedQuery, 1)
		require.NoError(t, err)
		defer res.Close()

		for res.Next() {
			var id int
			err := res.Scan(&id)
			require.NoError(t, err)
			require.Equal(t, 1, id)
		}

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return correct  QueryRowx result", func(t *testing.T) {
		expectedQuery := `SELECT * FROM "users" WHERE id = ?`

		mock.ExpectPrepare(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = ?`)).
			ExpectQuery().
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		res, err := sqlxx.QueryRowxContext(ctx, expectedQuery, 1)
		require.NoError(t, err)

		var id int
		err = res.Scan(&id)
		require.NoError(t, err)
		require.Equal(t, 1, id)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return correct  ExecContext result", func(t *testing.T) {
		expectedQuery := `INSERT INTO "users"(id) VALUES (?)`

		mock.ExpectPrepare(regexp.QuoteMeta(expectedQuery)).
			ExpectExec().
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		res, err := sqlxx.ExecContext(ctx, expectedQuery, 1)
		require.NoError(t, err)

		rowAffected, err := res.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), rowAffected)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
