package bank_account_test

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/bank_accounts"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/bank_account"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/rand"
	"testing"
)

func Test_service_GetAll(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	ctx := context.Background()

	mockBankAccouncRepository := bank_accounts.NewMockRepository(mock)

	s := bank_account.NewService(bank_account.NewServiceOpts{
		BankAccountRepository: mockBankAccouncRepository,
	})

	t.Run("should be return correct", func(t *testing.T) {
		expectedTotalData := int64(2)
		expectedPagination := pagination.PaginationInput{
			Page:     1,
			PageSize: 2,
		}
		expectedInput := bank_account.GetAllInput{
			ConsumerID: rand.Int63(),
			Pagination: expectedPagination,
		}
		expectedOutput := bank_account.GetAllOutput{
			Pagination: pagination.CreatePaginationOutput(expectedPagination, expectedTotalData),
			Items: []bank_account.GetAllOutputItem{
				{
					ID:                rand.Int63(),
					ConsumerID:        rand.Int63(),
					Name:              faker.Name(),
					AccountNumber:     faker.Name(),
					AccountHolderName: faker.Name(),
				},
				{
					ID:                rand.Int63(),
					ConsumerID:        rand.Int63(),
					Name:              faker.Name(),
					AccountNumber:     faker.Name(),
					AccountHolderName: faker.Name(),
				},
			},
		}

		mockBankAccouncRepository.EXPECT().
			GetAll(ctx, bank_accounts.GetAllInput{
				ConsumerID: null.IntFrom(expectedInput.ConsumerID),
				Pagination: expectedPagination,
			}).
			Return(bank_accounts.GetAllOutput{
				Pagination: pagination.CreatePaginationOutput(expectedPagination, expectedTotalData),
				Items: []bank_accounts.GetAllOutputItem{
					{
						ID:                expectedOutput.Items[0].ID,
						ConsumerID:        expectedOutput.Items[0].ConsumerID,
						Name:              expectedOutput.Items[0].Name,
						AccountNumber:     expectedOutput.Items[0].AccountNumber,
						AccountHolderName: expectedOutput.Items[0].AccountHolderName,
					},
					{
						ID:                expectedOutput.Items[1].ID,
						ConsumerID:        expectedOutput.Items[1].ConsumerID,
						Name:              expectedOutput.Items[1].Name,
						AccountNumber:     expectedOutput.Items[1].AccountNumber,
						AccountHolderName: expectedOutput.Items[1].AccountHolderName,
					},
				},
			}, nil)

		output, err := s.GetAll(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})
}
