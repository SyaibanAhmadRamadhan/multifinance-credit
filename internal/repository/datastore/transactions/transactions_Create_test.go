package transactions_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/transactions"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func Test_repository_Creates(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewSqlxWrapper(sqlxDB)

	r := transactions.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := transactions.CreateInput{
			LimitID:         rand.Int63(),
			ConsumerID:      rand.Int63(),
			ContractNumber:  rand.Int63(),
			Amount:          rand.Float64(),
			TransactionDate: time.Time{},
			Status:          "ACTIVE",
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`INSERT INTO transactions (limit_id,consumer_id,contract_number,amount,transaction_date,status,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?)`,
		)).ExpectExec().WithArgs(
			expectedInput.LimitID,
			expectedInput.ConsumerID,
			expectedInput.ContractNumber,
			expectedInput.Amount,
			expectedInput.TransactionDate,
			expectedInput.Status,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))

		output, err := r.Create(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, int64(1), output.ID)
	})
}
