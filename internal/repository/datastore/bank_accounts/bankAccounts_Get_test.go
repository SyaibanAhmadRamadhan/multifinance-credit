package bank_accounts_test

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/bank_accounts"
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

	r := bank_accounts.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := bank_accounts.GetInput{
			ID:            null.IntFrom(rand.Int63()),
			AccountNumber: null.StringFrom(faker.Name()),
		}

		expectedOutput := bank_accounts.GetOutput{
			ID:                expectedInput.ID.Int64,
			ConsumerID:        rand.Int63(),
			Name:              faker.Name(),
			AccountNumber:     expectedInput.AccountNumber.String,
			AccountHolderName: faker.Name(),
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT id, consumer_id, name, account_number, account_holder_name FROM bank_accounts WHERE account_number = ? AND id = ?`,
		)).ExpectQuery().WithArgs(expectedInput.AccountNumber.String, expectedInput.ID.Int64).
			WillReturnRows(sqlmock.NewRows([]string{"id", "consumer_id", "name", "account_number", "account_holder_name"}).
				AddRow(
					expectedOutput.ID,
					expectedOutput.ConsumerID,
					expectedOutput.Name,
					expectedOutput.AccountNumber,
					expectedOutput.AccountHolderName,
				))

		output, err := r.Get(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("should be return error no rows", func(t *testing.T) {
		expectedInput := bank_accounts.GetInput{
			ID:            null.IntFrom(rand.Int63()),
			AccountNumber: null.StringFrom(faker.Name()),
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT id, consumer_id, name, account_number, account_holder_name FROM bank_accounts WHERE account_number = ? AND id = ?`,
		)).ExpectQuery().WithArgs(expectedInput.AccountNumber.String, expectedInput.ID.Int64).
			WillReturnError(sql.ErrNoRows)

		output, err := r.Get(ctx, expectedInput)
		require.ErrorIs(t, err, datastore.ErrRecordNotFound)
		require.Empty(t, output)
	})
}
