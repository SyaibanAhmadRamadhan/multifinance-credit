package products_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/products"
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

	r := products.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedMerchantID := rand.Int63()
		expectedInput := products.CreatesInput{
			MerchantID: expectedMerchantID,
			Items: []products.CreatesInputItem{
				{
					Name:  faker.Name(),
					Image: faker.Name(),
					Qty:   rand.Int63(),
					Price: rand.Float64(),
				},
				{
					Name:  faker.Name(),
					Image: faker.Name(),
					Qty:   rand.Int63(),
					Price: rand.Float64(),
				},
			},
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`INSERT INTO products (merchant_id,image,name,qty,price) VALUES (?,?,?,?,?),(?,?,?,?,?)`,
		)).ExpectExec().WithArgs(
			expectedInput.MerchantID,
			expectedInput.Items[0].Image,
			expectedInput.Items[0].Name,
			expectedInput.Items[0].Qty,
			expectedInput.Items[0].Price,
			expectedInput.MerchantID,
			expectedInput.Items[1].Image,
			expectedInput.Items[1].Name,
			expectedInput.Items[1].Qty,
			expectedInput.Items[1].Price,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err = r.Creates(ctx, expectedInput)
		require.NoError(t, err)
	})
}
