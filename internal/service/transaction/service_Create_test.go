package transaction_test

import (
	"context"
	"database/sql"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/installments"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/limits"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/products"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/transaction_items"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/transactions"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/service/transaction"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	extra "github.com/oxyno-zeta/gomock-extra-matcher"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/rand"
	"testing"
	"time"
)

func Test_service_Create(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mockDBTx := db.NewMockSqlxTransaction(mock)
	mockLimitRepository := limits.NewMockRepository(mock)
	mockProductRepository := products.NewMockRepository(mock)
	mockTransactionRepository := transactions.NewMockRepository(mock)
	mockTransactionItemRepository := transaction_items.NewMockRepository(mock)
	mockInstallmentRepository := installments.NewMockRepository(mock)

	s := transaction.NewService(transaction.NewServiceOpts{
		ProductRepository:         mockProductRepository,
		TransactionRepository:     mockTransactionRepository,
		TransactionItemRepository: mockTransactionItemRepository,
		InstallmentRepository:     mockInstallmentRepository,
		LimitRepository:           mockLimitRepository,
		DBTx:                      mockDBTx,
	})

	ctx := context.Background()

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := transaction.CreateInput{
			Products: []transaction.CreateInputProductItem{
				{
					ID:  rand.Int63(),
					Qty: 3,
				},
				{
					ID:  rand.Int63(),
					Qty: 2,
				},
			},
			LimitID:    rand.Int63(),
			ConsumerID: rand.Int63(),
		}
		expectedLimitOutputRepository := limits.GetOutput{
			ID:              expectedInput.LimitID,
			ConsumerID:      expectedInput.ConsumerID,
			Tenor:           3,
			Amount:          10_000_000,
			RemainingAmount: 10_000_000,
		}
		expectedProductsOutputRepository := products.GetAllOutput{
			Items: []products.GetAllOutputItem{
				{
					ID:         expectedInput.Products[0].ID,
					MerchantID: rand.Int63(),
					Image:      faker.UUIDDigit(),
					Name:       faker.Name(),
					Qty:        5,
					Price:      100_000,
				},
				{
					ID:         expectedInput.Products[1].ID,
					MerchantID: rand.Int63(),
					Image:      faker.UUIDDigit(),
					Name:       faker.Name(),
					Qty:        2,
					Price:      1_000_000,
				},
			},
		}
		expectedTotalAmount := 3*100_000 + 2*1_000_000
		expectedRemainingLimit := expectedLimitOutputRepository.RemainingAmount - float64(expectedTotalAmount)
		expectedTransactionID := rand.Int63()
		timeNow := time.Now().UTC()
		timeNow = time.Date(timeNow.Year(), timeNow.Month(), 1, 0, 0, 0, 0, time.UTC)
		calculateOutput := util.CalculateMonthlyInstallments(util.CalculateMonthlyInstallmentsInput{
			Principal:          float64(expectedTotalAmount),
			AnnualInterestRate: 2.0,
			Tenor:              expectedLimitOutputRepository.Tenor,
		})
		expectedInstallmentInputItem := []installments.CreatesInputItem{
			{
				Amount:  calculateOutput.MonthlyInstallments,
				DueDate: timeNow.AddDate(0, 1, 0),
				Status:  string(primitive.InstallmentStatusUnPaid),
			},
			{
				Amount:  calculateOutput.MonthlyInstallments,
				DueDate: timeNow.AddDate(0, 2, 0),
				Status:  string(primitive.InstallmentStatusUnPaid),
			},
			{
				Amount:  calculateOutput.MonthlyInstallments,
				DueDate: timeNow.AddDate(0, 3, 0),
				Status:  string(primitive.InstallmentStatusUnPaid),
			},
		}

		mockDBTx.EXPECT().DoTransaction(
			ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false}, gomock.Any(),
		).DoAndReturn(func(ctx context.Context, tx *sql.TxOptions, fn func(tx db.Rdbms) error) error {
			mockTx := db.NewRdbms(nil)

			mockLimitRepository.EXPECT().
				Get(ctx, limits.GetInput{
					Tx:         mockTx,
					Locking:    db.LockingUpdate,
					ID:         null.IntFrom(expectedInput.LimitID),
					ConsumerID: null.IntFrom(expectedInput.ConsumerID),
				}).Return(expectedLimitOutputRepository, nil)

			mockProductRepository.EXPECT().
				GetAll(ctx, products.GetAllInput{
					Locking:     db.LockingUpdate,
					Transaction: mockTx,
					IDs:         []int64{expectedInput.Products[0].ID, expectedInput.Products[1].ID},
					Pagination: pagination.PaginationInput{
						Page:     1,
						PageSize: 2,
					},
				}).Return(expectedProductsOutputRepository, nil)

			createTransactionInput := extra.StructMatcher().
				Field("Transaction", mockTx).
				Field("LimitID", expectedInput.LimitID).
				Field("ConsumerID", expectedInput.ConsumerID).
				Field("ContractNumber", gomock.Any()).
				Field("Amount", float64(expectedTotalAmount)).
				Field("TransactionDate", gomock.Any()).
				Field("Status", string(primitive.TransactionStatusActive))
			mockTransactionRepository.EXPECT().
				Create(ctx, createTransactionInput).
				Return(transactions.CreateOutput{
					ID: expectedTransactionID,
				}, nil)

			mockTransactionItemRepository.EXPECT().
				Creates(ctx, transaction_items.CreatesInput{
					Transaction:   mockTx,
					TransactionID: expectedTransactionID,
					Items: []transaction_items.CreatesItemInput{
						{
							MerchantID: expectedProductsOutputRepository.Items[0].MerchantID,
							Name:       expectedProductsOutputRepository.Items[0].Name,
							Image:      expectedProductsOutputRepository.Items[0].Image,
							Qty:        expectedProductsOutputRepository.Items[0].Qty,
							UnitPrice:  expectedProductsOutputRepository.Items[0].Price,
							Amount:     float64(expectedInput.Products[0].Qty) * expectedProductsOutputRepository.Items[0].Price,
						},
						{
							MerchantID: expectedProductsOutputRepository.Items[1].MerchantID,
							Name:       expectedProductsOutputRepository.Items[1].Name,
							Image:      expectedProductsOutputRepository.Items[1].Image,
							Qty:        expectedProductsOutputRepository.Items[1].Qty,
							UnitPrice:  expectedProductsOutputRepository.Items[1].Price,
							Amount:     float64(expectedInput.Products[1].Qty) * expectedProductsOutputRepository.Items[1].Price,
						},
					},
				}).Return(nil)

			mockProductRepository.EXPECT().
				Updates(ctx, products.UpdatesInput{
					Transaction: mockTx,
					Items: []products.UpdatesInputItem{
						{
							ID:  expectedInput.Products[0].ID,
							Qty: expectedProductsOutputRepository.Items[0].Qty - expectedInput.Products[0].Qty,
						},
						{
							ID:  expectedInput.Products[1].ID,
							Qty: expectedProductsOutputRepository.Items[1].Qty - expectedInput.Products[1].Qty,
						},
					},
				}).Return(nil)

			createInstallmentInput := extra.StructMatcher().
				Field("Transaction", mockTx).
				Field("LimitID", expectedInput.LimitID).
				Field("ContractNumber", gomock.Any()).
				Field("Items", gomock.InAnyOrder(expectedInstallmentInputItem))
			mockInstallmentRepository.EXPECT().
				Creates(ctx, createInstallmentInput).
				Return(nil)

			mockLimitRepository.EXPECT().
				Update(ctx, limits.UpdateInput{
					Transaction:     mockTx,
					ID:              expectedInput.LimitID,
					RemainingAmount: expectedRemainingLimit,
				}).Return(nil)

			return fn(mockTx)
		})

		output, err := s.Create(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedTransactionID, output.ID)
	})
}
