package users_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/users"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
)

func Test_repository_CheckExisting(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewRdbms(sqlxDB)

	r := users.NewRepository(sqlxx)

	t.Run("should be return correct with existing data", func(t *testing.T) {
		expectedInput := users.CheckExistingInput{
			ByID:    null.IntFrom(rand.Int63()),
			ByEmail: null.StringFrom(faker.Email()),
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT EXISTS( SELECT 1 FROM users WHERE id = ? AND email = ? )`)).
			WithArgs(expectedInput.ByID.Int64, expectedInput.ByEmail.String).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		output, err := r.CheckExisting(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, true, output.Existing)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should be return correct without existing data", func(t *testing.T) {
		expectedInput := users.CheckExistingInput{
			ByID:    null.IntFrom(rand.Int63()),
			ByEmail: null.StringFrom(faker.Email()),
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT EXISTS( SELECT 1 FROM users WHERE id = ? AND email = ? )`)).
			WithArgs(expectedInput.ByID.Int64, expectedInput.ByEmail.String).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		output, err := r.CheckExisting(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, false, output.Existing)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
