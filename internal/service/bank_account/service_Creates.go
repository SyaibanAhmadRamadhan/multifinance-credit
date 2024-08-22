package bank_account

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/bank_accounts"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (s *service) Creates(ctx context.Context, input CreatesInput) (err error) {
	const batch = 10

	itemBatches := util.SplitDataIntoBatch(input.Items, batch)
	if len(itemBatches) <= 0 {
		return
	}

	err = s.dbTx.DoTransaction(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	}, func(tx *db.SqlxWrapper) error {

		for i, items := range itemBatches {
			createsBankAccountItems := make([]bank_accounts.CreatesInputItem, 0)

			for _, item := range items {
				createsBankAccountItems = append(createsBankAccountItems, bank_accounts.CreatesInputItem{
					ConsumerID:        input.ConsumerID,
					Name:              item.Name,
					AccountNumber:     item.AccountNumber,
					AccountHolderName: item.AccountHolderName,
				})
			}

			err = s.bankAccountRepository.Creates(ctx, bank_accounts.CreatesInput{
				Transaction: tx,
				Items:       createsBankAccountItems,
			})

			if err != nil {
				return tracer.Error(fmt.Errorf("errors in batch: %d. %w", i, err))
			}
		}

		return nil
	})
	if err != nil {
		return tracer.Error(err)
	}

	return
}
