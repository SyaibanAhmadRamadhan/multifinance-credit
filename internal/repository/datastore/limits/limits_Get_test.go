package limits_test

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/limits"
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

	sqlxx := db.NewRdbms(sqlxDB)

	r := limits.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := limits.GetInput{
			ID: null.IntFrom(rand.Int63()),
		}

		expectedOutput := limits.GetOutput{
			ID:              expectedInput.ID.Int64,
			ConsumerID:      rand.Int63(),
			Tenor:           rand.Int31(),
			Amount:          rand.Float64(),
			RemainingAmount: rand.Float64(),
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, consumer_id, tenor, amount, remaining_amount FROM limits WHERE id = ?`,
		)).WithArgs(expectedInput.ID.Int64).
			WillReturnRows(sqlmock.NewRows([]string{"id", "consumer_id", "tenor", "amount", "remaining_amount"}).
				AddRow(
					expectedOutput.ID,
					expectedOutput.ConsumerID,
					expectedOutput.Tenor,
					expectedOutput.Amount,
					expectedOutput.RemainingAmount,
				))

		output, err := r.Get(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("should be return error no rows", func(t *testing.T) {
		expectedInput := limits.GetInput{
			ID: null.IntFrom(rand.Int63()),
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, consumer_id, tenor, amount, remaining_amount FROM limits WHERE id = ?`,
		)).WithArgs(expectedInput.ID.Int64).
			WillReturnError(sql.ErrNoRows)

		output, err := r.Get(ctx, expectedInput)
		require.ErrorIs(t, err, datastore.ErrRecordNotFound)
		require.Empty(t, output)
	})
}
