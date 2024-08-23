package installments_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/installments"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func Test_repository_GetAll(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewSqlxWrapper(sqlxDB)

	r := installments.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedTotalData := int64(1)
		expectedPaginationInput := pagination.PaginationInput{
			Page:     1,
			PageSize: 2,
		}

		expectedInput := installments.GetAllInput{
			Pagination:     expectedPaginationInput,
			ContractNumber: rand.Int63(),
		}

		expectedOutput := installments.GetAllOutput{
			Pagination: pagination.CreatePaginationOutput(expectedPaginationInput, expectedTotalData),
			Items: []installments.GetAllOutputItem{
				{
					ID:              rand.Int63(),
					LimitID:         rand.Int63(),
					PaymentMethodID: nil,
					ContractNumber:  expectedInput.ContractNumber,
					Amount:          rand.Float64(),
					DueDate:         time.Now().UTC(),
					PaymentDate:     null.TimeFrom(time.Now().UTC()).Ptr(),
					Status:          "PAID",
				},
				{
					ID:              rand.Int63(),
					LimitID:         rand.Int63(),
					PaymentMethodID: nil,
					ContractNumber:  expectedInput.ContractNumber,
					Amount:          rand.Float64(),
					DueDate:         time.Now().UTC(),
					PaymentDate:     null.TimeFrom(time.Now().UTC()).Ptr(),
					Status:          "PAID",
				},
			},
		}

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT COUNT(*) FROM installments WHERE contract_number = ?`,
		)).ExpectQuery().WithArgs(expectedInput.ContractNumber).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedTotalData))

		mock.ExpectPrepare(regexp.QuoteMeta(
			`SELECT id, limit_id, contract_number, amount, due_date, payment_date, status FROM installments WHERE contract_number = ? ORDER BY due_date ASC, id DESC LIMIT 2 OFFSET 0`,
		)).ExpectQuery().WithArgs(expectedInput.ContractNumber).
			WillReturnRows(sqlmock.NewRows([]string{"id", "limit_id", "contract_number", "amount", "due_date", "payment_date", "status"}).
				AddRow(
					expectedOutput.Items[0].ID,
					expectedOutput.Items[0].LimitID,
					expectedOutput.Items[0].ContractNumber,
					expectedOutput.Items[0].Amount,
					expectedOutput.Items[0].DueDate,
					expectedOutput.Items[0].PaymentDate,
					expectedOutput.Items[0].Status,
				).
				AddRow(
					expectedOutput.Items[1].ID,
					expectedOutput.Items[1].LimitID,
					expectedOutput.Items[1].ContractNumber,
					expectedOutput.Items[1].Amount,
					expectedOutput.Items[1].DueDate,
					expectedOutput.Items[1].PaymentDate,
					expectedOutput.Items[1].Status,
				))

		output, err := r.GetAll(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})
}
