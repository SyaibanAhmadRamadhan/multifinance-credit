package products_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/products"
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

	sqlxx := db.NewSqlxWrapper(sqlxDB)

	r := products.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedTotalData := int64(1)
		expectedPaginationInput := pagination.PaginationInput{
			Page:     1,
			PageSize: 2,
		}

		expectedInput := products.GetAllInput{
			Pagination: expectedPaginationInput,
			MerchantID: null.IntFrom(rand.Int63()),
			IDs:        []int64{1, 2},
		}

		expectedOutput := products.GetAllOutput{
			Pagination: pagination.CreatePaginationOutput(expectedPaginationInput, expectedTotalData),
			Items: []products.GetAllOutputItem{
				{
					ID:         rand.Int63(),
					MerchantID: expectedInput.MerchantID.Int64,
					Image:      faker.Name(),
					Name:       faker.Name(),
					Qty:        rand.Int63(),
					Price:      rand.Float64(),
				},
				{
					ID:         rand.Int63(),
					MerchantID: expectedInput.MerchantID.Int64,
					Image:      faker.Name(),
					Name:       faker.Name(),
					Qty:        rand.Int63(),
					Price:      rand.Float64(),
				},
			},
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT COUNT(*) FROM products WHERE merchant_id = ? AND id IN (?,?)`,
		)).ExpectQuery().WithArgs(expectedInput.MerchantID.Int64, expectedInput.IDs[0], expectedInput.IDs[1]).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedTotalData))

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT id, merchant_id, name, image, qty, price FROM products WHERE merchant_id = ? AND id IN (?,?) LIMIT 2 OFFSET 0`,
		)).ExpectQuery().WithArgs(expectedInput.MerchantID.Int64, expectedInput.IDs[0], expectedInput.IDs[1]).
			WillReturnRows(sqlmock.NewRows([]string{"id", "merchant_id", "name", "image", "qty", "price"}).
				AddRow(
					expectedOutput.Items[0].ID,
					expectedOutput.Items[0].MerchantID,
					expectedOutput.Items[0].Name,
					expectedOutput.Items[0].Image,
					expectedOutput.Items[0].Qty,
					expectedOutput.Items[0].Price,
				).
				AddRow(
					expectedOutput.Items[1].ID,
					expectedOutput.Items[1].MerchantID,
					expectedOutput.Items[1].Name,
					expectedOutput.Items[1].Image,
					expectedOutput.Items[1].Qty,
					expectedOutput.Items[1].Price,
				))

		output, err := r.GetAll(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})
}
