package bank_account

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/bank_accounts"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/guregu/null/v5"
)

func (s *service) GetAll(ctx context.Context, input GetAllInput) (output GetAllOutput, err error) {
	bankAccountsOutput, err := s.bankAccountRepository.GetAll(ctx, bank_accounts.GetAllInput{
		ConsumerID: null.IntFrom(input.ConsumerID),
		Pagination: input.Pagination,
	})
	if err != nil {
		return output, tracer.Error(err)
	}

	output = GetAllOutput{
		Pagination: bankAccountsOutput.Pagination,
		Items:      make([]GetAllOutputItem, 0),
	}

	for _, item := range bankAccountsOutput.Items {
		output.Items = append(output.Items, GetAllOutputItem{
			ID:                item.ID,
			ConsumerID:        item.ConsumerID,
			Name:              item.Name,
			AccountNumber:     item.AccountNumber,
			AccountHolderName: item.AccountHolderName,
		})
	}
	return
}
