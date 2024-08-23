package products_test

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/products"
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

	r := products.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := products.GetInput{
			ID: null.IntFrom(rand.Int63()),
		}

		expectedOutput := products.GetOutput{
			ID:         expectedInput.ID.Int64,
			MerchantID: rand.Int63(),
			Image:      faker.Name(),
			Name:       faker.Name(),
			Qty:        rand.Int63(),
			Price:      rand.Float64(),
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT id, merchant_id, name, image, qty, price FROM products WHERE id = ?`,
		)).ExpectQuery().WithArgs(expectedInput.ID.Int64).
			WillReturnRows(sqlmock.NewRows([]string{"id", "merchant_id", "name", "image", "qty", "price"}).
				AddRow(
					expectedOutput.ID,
					expectedOutput.MerchantID,
					expectedOutput.Name,
					expectedOutput.Image,
					expectedOutput.Qty,
					expectedOutput.Price,
				))

		output, err := r.Get(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("should be return error no rows", func(t *testing.T) {
		expectedInput := products.GetInput{
			ID: null.IntFrom(rand.Int63()),
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT id, merchant_id, name, image, qty, price FROM products WHERE id = ?`,
		)).ExpectQuery().WithArgs(expectedInput.ID.Int64).
			WillReturnError(sql.ErrNoRows)

		output, err := r.Get(ctx, expectedInput)
		require.ErrorIs(t, err, datastore.ErrRecordNotFound)
		require.Empty(t, output)
	})
}
