package transaction_items_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/transaction_items"
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

	sqlxx := db.NewSqlxWrapper(sqlxDB)

	r := transaction_items.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := transaction_items.CreatesInput{
			TransactionID: rand.Int63(),
			Items: []transaction_items.CreatesItemInput{
				{
					MerchantID: rand.Int63(),
					Name:       faker.Name(),
					Image:      faker.Name(),
					Qty:        rand.Int63(),
					UnitPrice:  rand.Float64(),
					Amount:     rand.Float64(),
				},
			},
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`INSERT INTO transaction_items (transaction_id,merchant_id,name,image,qty,unit_price,amount) VALUES (?,?,?,?,?,?,?)`,
		)).ExpectExec().WithArgs(
			expectedInput.TransactionID,
			expectedInput.Items[0].MerchantID,
			expectedInput.Items[0].Name,
			expectedInput.Items[0].Image,
			expectedInput.Items[0].Qty,
			expectedInput.Items[0].UnitPrice,
			expectedInput.Items[0].Amount,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err = r.Creates(ctx, expectedInput)
		require.NoError(t, err)
	})
}
