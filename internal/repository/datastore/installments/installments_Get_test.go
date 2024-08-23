package installments_test

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/installments"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func Test_repository_Get(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewSqlxWrapper(sqlxDB)

	r := installments.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := installments.GetInput{
			ID: null.IntFrom(rand.Int63()),
		}

		expectedOutput := installments.GetOutput{
			ID:              expectedInput.ID.Int64,
			LimitID:         rand.Int63(),
			PaymentMethodID: nil,
			ContractNumber:  rand.Int63(),
			Amount:          rand.Float64(),
			DueDate:         time.Now().UTC(),
			PaymentDate:     nil,
			Status:          "UNPAID",
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT id, limit_id, contract_number, amount, due_date, payment_date, status FROM installments WHERE id = ?`,
		)).ExpectQuery().WithArgs(expectedInput.ID.Int64).
			WillReturnRows(sqlmock.NewRows([]string{"id", "limit_id", "contract_number", "amount", "due_date", "payment_date", "status"}).
				AddRow(
					expectedOutput.ID,
					expectedOutput.LimitID,
					expectedOutput.ContractNumber,
					expectedOutput.Amount,
					expectedOutput.DueDate,
					expectedOutput.PaymentDate,
					expectedOutput.Status,
				))

		output, err := r.Get(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("should be return error no rows", func(t *testing.T) {
		expectedInput := installments.GetInput{
			ID: null.IntFrom(rand.Int63()),
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT id, limit_id, contract_number, amount, due_date, payment_date, status FROM installments WHERE id = ?`,
		)).ExpectQuery().WithArgs(expectedInput.ID.Int64).
			WillReturnError(sql.ErrNoRows)

		output, err := r.Get(ctx, expectedInput)
		require.ErrorIs(t, err, datastore.ErrRecordNotFound)
		require.Empty(t, output)
	})
}
