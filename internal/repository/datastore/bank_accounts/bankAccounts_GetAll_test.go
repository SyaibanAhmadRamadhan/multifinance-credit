package bank_accounts_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/bank_accounts"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
)

func Test_repository_GetAll(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewRdbms(sqlxDB)

	r := bank_accounts.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedTotalData := int64(1)
		expectedPaginationInput := pagination.PaginationInput{
			Page:     1,
			PageSize: 2,
		}

		expectedInput := bank_accounts.GetAllInput{
			Pagination: expectedPaginationInput,
			ConsumerID: null.IntFrom(rand.Int63()),
		}

		expectedOutput := bank_accounts.GetAllOutput{
			Pagination: pagination.CreatePaginationOutput(expectedPaginationInput, expectedTotalData),
			Items: []bank_accounts.GetAllOutputItem{
				{
					ID:                rand.Int63(),
					ConsumerID:        expectedInput.ConsumerID.Int64,
					Name:              faker.Name(),
					AccountNumber:     faker.Name(),
					AccountHolderName: faker.Name(),
				},
				{
					ID:                rand.Int63(),
					ConsumerID:        expectedInput.ConsumerID.Int64,
					Name:              faker.Name(),
					AccountNumber:     faker.Name(),
					AccountHolderName: faker.Name(),
				},
			},
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT COUNT(*) FROM bank_accounts WHERE consumer_id = ?`,
		)).WithArgs(expectedInput.ConsumerID.Int64).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedTotalData))

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, consumer_id, name, account_number, account_holder_name FROM bank_accounts WHERE consumer_id = ? LIMIT 2`,
		)).WithArgs(expectedInput.ConsumerID.Int64).
			WillReturnRows(sqlmock.NewRows([]string{"id", "consumer_id", "name", "account_number", "account_holder_name"}).
				AddRow(
					expectedOutput.Items[0].ID,
					expectedOutput.Items[0].ConsumerID,
					expectedOutput.Items[0].Name,
					expectedOutput.Items[0].AccountNumber,
					expectedOutput.Items[0].AccountHolderName,
				).
				AddRow(
					expectedOutput.Items[1].ID,
					expectedOutput.Items[1].ConsumerID,
					expectedOutput.Items[1].Name,
					expectedOutput.Items[1].AccountNumber,
					expectedOutput.Items[1].AccountHolderName,
				))

		output, err := r.GetAll(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})
}
