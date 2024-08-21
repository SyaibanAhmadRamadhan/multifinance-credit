package users_test

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/users"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
)

func Test_repository_Get(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewSqlxWrapper(sqlxDB)

	r := users.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := users.GetInput{
			ID:    null.IntFrom(rand.Int63()),
			Email: null.StringFrom(faker.Email()),
		}
		expectedOutput := users.GetOutput{
			ID:       expectedInput.ID.Int64,
			Email:    expectedInput.Email.String,
			Password: faker.Password(),
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT id, email, password FROM users WHERE email = ? AND id = ?`,
		)).ExpectQuery().WithArgs(expectedInput.Email.String, expectedInput.ID.Int64).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "email", "password",
			}).AddRow(
				expectedOutput.ID, expectedOutput.Email, expectedOutput.Password,
			))

		output, err := r.Get(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should be return error no record found", func(t *testing.T) {
		expectedInput := users.GetInput{
			ID:    null.IntFrom(rand.Int63()),
			Email: null.StringFrom(faker.Email()),
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT id, email, password FROM users WHERE email = ? AND id = ?`,
		)).ExpectQuery().WithArgs(expectedInput.Email.String, expectedInput.ID.Int64).
			WillReturnError(sql.ErrNoRows)

		output, err := r.Get(ctx, expectedInput)
		require.Error(t, err, datastore.ErrRecordNotFound)
		require.Empty(t, output)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
