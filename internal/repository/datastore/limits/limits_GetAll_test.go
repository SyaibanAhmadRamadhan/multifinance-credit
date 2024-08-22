package limits_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/limits"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
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

	r := limits.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedTotalData := int64(1)
		expectedPaginationInput := pagination.PaginationInput{
			Page:     1,
			PageSize: 2,
		}

		expectedInput := limits.GetAllInput{
			Pagination: expectedPaginationInput,
			ConsumerID: null.IntFrom(rand.Int63()),
		}

		expectedOutput := limits.GetAllOutput{
			Pagination: pagination.CreatePaginationOutput(expectedPaginationInput, expectedTotalData),
			Items: []limits.GetAllOutputItem{
				{
					ID:         rand.Int63(),
					ConsumerID: expectedInput.ConsumerID.Int64,
					Tenor:      rand.Int31(),
					Amount:     rand.Float64(),
				},
				{
					ID:         rand.Int63(),
					ConsumerID: expectedInput.ConsumerID.Int64,
					Tenor:      rand.Int31(),
					Amount:     rand.Float64(),
				},
			},
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT COUNT(*) FROM limits WHERE consumer_id = ?`,
		)).ExpectQuery().WithArgs(expectedInput.ConsumerID.Int64).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedTotalData))

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT id, consumer_id, tenor, amount FROM limits WHERE consumer_id = ? LIMIT 2`,
		)).ExpectQuery().WithArgs(expectedInput.ConsumerID.Int64).
			WillReturnRows(sqlmock.NewRows([]string{"id", "consumer_id", "tenor", "amount"}).
				AddRow(
					expectedOutput.Items[0].ID,
					expectedOutput.Items[0].ConsumerID,
					expectedOutput.Items[0].Tenor,
					expectedOutput.Items[0].Amount,
				).
				AddRow(
					expectedOutput.Items[1].ID,
					expectedOutput.Items[1].ConsumerID,
					expectedOutput.Items[1].Tenor,
					expectedOutput.Items[1].Amount,
				))

		output, err := r.GetAll(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})
}
