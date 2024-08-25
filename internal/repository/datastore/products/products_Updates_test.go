package products_test

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/products"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
)

func Test_getItemIDs(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewRdbms(sqlxDB)

	r := products.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := products.UpdatesInput{
			Items: []products.UpdatesInputItem{
				{
					ID:  rand.Int63(),
					Qty: rand.Int63(),
				},
				{
					ID:  rand.Int63(),
					Qty: rand.Int63(),
				},
			},
		}
		mock.ExpectExec(regexp.QuoteMeta(
			fmt.Sprintf(`UPDATE products SET qty = CASE id WHEN %d THEN %d WHEN %d THEN %d END WHERE id IN (?,?)`,
				expectedInput.Items[0].ID, expectedInput.Items[0].Qty, expectedInput.Items[1].ID, expectedInput.Items[1].Qty),
		)).WithArgs(expectedInput.Items[0].ID, expectedInput.Items[1].ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = r.Updates(ctx, expectedInput)
		require.NoError(t, err)
	})
}
