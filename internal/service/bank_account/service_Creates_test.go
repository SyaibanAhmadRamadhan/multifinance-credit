package bank_account_test

import (
	"context"
	"database/sql"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/bank_accounts"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/bank_account"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/rand"
	"testing"
)

func Test_service_Creates(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()

	mockBankAccouncRepository := bank_accounts.NewMockRepository(mock)
	mockDBTx := db.NewMockSqlxTransaction(mock)

	s := bank_account.NewService(bank_account.NewServiceOpts{
		BankAccountRepository: mockBankAccouncRepository,
		DBTx:                  mockDBTx,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := bank_account.CreatesInput{
			ConsumerID: rand.Int63(),
			Items: []bank_account.CreatesInputItem{
				{
					Name:              faker.Name(),
					AccountNumber:     faker.Name(),
					AccountHolderName: faker.Name(),
				},
				{
					Name:              faker.Name(),
					AccountNumber:     faker.Name(),
					AccountHolderName: faker.Name(),
				},
			},
		}

		mockDBTx.EXPECT().
			DoTransaction(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false},
				gomock.Any()).
			DoAndReturn(func(ctx context.Context, tx *sql.TxOptions, fn func(tx *db.SqlxWrapper) error) error {
				mockSqlxWrapper := &db.SqlxWrapper{}
				mockBankAccouncRepository.EXPECT().
					Creates(ctx, bank_accounts.CreatesInput{
						Transaction: mockSqlxWrapper,
						Items: []bank_accounts.CreatesInputItem{
							{
								ConsumerID:        expectedInput.ConsumerID,
								Name:              expectedInput.Items[0].Name,
								AccountNumber:     expectedInput.Items[0].AccountNumber,
								AccountHolderName: expectedInput.Items[0].AccountHolderName,
							},
							{
								ConsumerID:        expectedInput.ConsumerID,
								Name:              expectedInput.Items[1].Name,
								AccountNumber:     expectedInput.Items[1].AccountNumber,
								AccountHolderName: expectedInput.Items[1].AccountHolderName,
							},
						},
					}).Return(nil)
				return fn(mockSqlxWrapper)
			}).Return(nil)

		err := s.Creates(ctx, expectedInput)
		require.NoError(t, err)
	})
}
