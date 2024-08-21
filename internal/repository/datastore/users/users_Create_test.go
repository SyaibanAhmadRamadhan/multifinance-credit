package users_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/users"
	"github.com/go-faker/faker/v4"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
)

func Test_repository_Create(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewSqlxWrapper(sqlxDB)

	r := users.NewRepository(sqlxx)

	t.Run("should be return correct without transaction", func(t *testing.T) {
		expectedInput := users.CreateInput{
			Transaction: nil,
			Email:       faker.Email(),
			Password:    faker.UUIDDigit(),
		}
		expectedID := rand.Int63()

		mock.ExpectPrepare(regexp.QuoteMeta(
			`INSERT INTO users (email,password,created_at) VALUES (?,?,?)`,
		)).ExpectExec().WithArgs(expectedInput.Email, expectedInput.Password, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(expectedID, 1))

		output, err := r.Create(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedID, output.ID)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}
