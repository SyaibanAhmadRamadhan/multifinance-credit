package bank_accounts_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/bank_accounts"
	"github.com/go-faker/faker/v4"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
)

func Test_repository_Creates(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewRdbms(sqlxDB)

	r := bank_accounts.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedConsumerID := rand.Int63()
		expectedInput := bank_accounts.CreatesInput{
			Items: []bank_accounts.CreatesInputItem{
				{
					ConsumerID:        expectedConsumerID,
					Name:              faker.Name(),
					AccountNumber:     faker.Name(),
					AccountHolderName: faker.Name(),
				},
				{
					ConsumerID:        expectedConsumerID,
					Name:              faker.Name(),
					AccountNumber:     faker.Name(),
					AccountHolderName: faker.Name(),
				},
			},
		}

		mock.ExpectExec(regexp.QuoteMeta(
			`INSERT INTO bank_accounts (consumer_id,name,account_number,account_holder_name) VALUES (?,?,?,?),(?,?,?,?)`,
		)).WithArgs(
			expectedInput.Items[0].ConsumerID,
			expectedInput.Items[0].Name,
			expectedInput.Items[0].AccountNumber,
			expectedInput.Items[0].AccountHolderName,
			expectedInput.Items[1].ConsumerID,
			expectedInput.Items[1].Name,
			expectedInput.Items[1].AccountNumber,
			expectedInput.Items[1].AccountHolderName,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err = r.Creates(ctx, expectedInput)
		require.NoError(t, err)
	})
}
