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

func Test_repository_Update(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := db.NewRdbms(sqlxDB)

	r := limits.NewRepository(sqlxx)

	t.Run("Should be return correct", func(t *testing.T) {
		expectedInput := limits.UpdateInput{
			ID:              rand.Int63(),
			RemainingAmount: rand.Float64(),
		}

		mock.ExpectExec(regexp.QuoteMeta(
			`UPDATE limits SET remaining_amount = ? WHERE id = ?`,
		)).WithArgs(expectedInput.RemainingAmount, expectedInput.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = r.Update(ctx, expectedInput)
		require.NoError(t, err)
	})
}
