package limits_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/limits"
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

	r := limits.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedConsumerID := rand.Int63()
		expectedInput := limits.CreatesInput{
			ConsumerID: expectedConsumerID,
			Items: []limits.CreatesInputItem{
				{

					Tenor:  rand.Int31(),
					Amount: rand.Float64(),
				},
				{
					Tenor:  rand.Int31(),
					Amount: rand.Float64(),
				},
			},
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`INSERT INTO limits (consumer_id,tenor,amount,remaining_amount) VALUES (?,?,?,?),(?,?,?,?)`,
		)).ExpectExec().WithArgs(
			expectedInput.ConsumerID,
			expectedInput.Items[0].Tenor,
			expectedInput.Items[0].Amount,
			expectedInput.Items[0].Amount,
			expectedInput.ConsumerID,
			expectedInput.Items[1].Tenor,
			expectedInput.Items[1].Amount,
			expectedInput.Items[1].Amount,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err = r.Creates(ctx, expectedInput)
		require.NoError(t, err)
	})
}
