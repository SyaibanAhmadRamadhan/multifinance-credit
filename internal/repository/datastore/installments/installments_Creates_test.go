package installments_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/installments"
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

	r := installments.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := installments.CreatesInput{
			LimitID:        rand.Int63(),
			ContractNumber: rand.Int63(),
			Items: []installments.CreatesInputItem{
				{
					Amount:  rand.Float64(),
					DueDate: time.Now().UTC(),
					Status:  "ACTIVE",
				},
				{
					Amount:  rand.Float64(),
					DueDate: time.Now().UTC(),
					Status:  "ACTIVE",
				},
			},
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`INSERT INTO installments (limit_id,contract_number,amount,due_date,status) VALUES (?,?,?,?,?),(?,?,?,?,?)`,
		)).ExpectExec().WithArgs(
			expectedInput.LimitID,
			expectedInput.ContractNumber,
			expectedInput.Items[0].Amount,
			expectedInput.Items[0].DueDate,
			expectedInput.Items[0].Status,
			expectedInput.LimitID,
			expectedInput.ContractNumber,
			expectedInput.Items[1].Amount,
			expectedInput.Items[1].DueDate,
			expectedInput.Items[1].Status,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err = r.Creates(ctx, expectedInput)
		require.NoError(t, err)
	})
}
