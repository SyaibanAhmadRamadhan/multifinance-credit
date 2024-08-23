package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/installments"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/limits"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/products"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/transaction_items"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/transactions"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/guregu/null/v5"
	"golang.org/x/sync/errgroup"
	"time"
)

func (s *service) Create(ctx context.Context, input CreateInput) (output CreateOutput, err error) {
	productID := make([]int64, 0)
	for _, product := range input.Products {
		productID = append(productID, product.ID)
	}

	err = s.dbTx.DoTransaction(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	}, func(tx *db.SqlxWrapper) (err error) {
		// GET DATA LIMIT
		limitOutput, err := s.limitRepository.Get(ctx, limits.GetInput{
			Tx:         tx,
			Locking:    db.LockingUpdate,
			ID:         null.IntFrom(input.LimitID),
			Tenor:      null.Int32{},
			ConsumerID: null.IntFrom(input.ConsumerID),
		})
		if err != nil {
			if errors.Is(err, datastore.ErrRecordNotFound) {
				err = errors.Join(err, ErrLimitNotFound)
			}
			fmt.Println(err)
			return tracer.Error(err)
		}

		// GET DATA PRODUCTS, FOR THE NEXT DEV YOU CAN GET DATA PER BATCH
		productsOutput, err := s.productRepository.GetAll(ctx, products.GetAllInput{
			Locking:     db.LockingUpdate,
			Transaction: tx,
			IDs:         productID,
			Pagination: pagination.PaginationInput{
				Page:     1,
				PageSize: int64(len(productID)),
			},
		})
		if err != nil {
			return tracer.Error(err)
		}

		// VALIDATE FOR PROCESSING TRANSACTION
		validateOutput, err := s.validateCreateTransaction(validateCreateInput{
			limit:       limitOutput,
			products:    productsOutput.Items,
			createInput: input,
		})
		if err != nil {
			return tracer.Error(err)
		}

		contractNumber := util.GenerateContractNumber()

		// USING ERR GROUP FOR TRIGGER ERROR WHEN ONE OF GOROUTINE RETURN ERROR
		var erg errgroup.Group

		// BACKGROUND PROCESS FOR INSERT DATA TO TRANSACTION TABLE AND TRANSACTION ITEM TABLE AND UPDATE QTY PRODUCT
		erg.Go(func() (err error) {
			output.ID, err = s.processInsertDataTransactionAndUpdateProduct(ctx, processInsertDataTransactionAndUpdateProductInput{
				tx:             tx,
				products:       productsOutput.Items,
				contractNumber: contractNumber,
				limitID:        limitOutput.ID,
				totalAmount:    validateOutput.totalAmount,
				consumerID:     input.ConsumerID,
				createInput:    input,
			})
			if err != nil {
				return tracer.Error(err)
			}
			return
		})

		// BACKGROUND PROCESS FOR INSERT DATA TO INSTALLMENT TABLE AND UPDATE REMAINING LIMIT
		erg.Go(func() (err error) {
			err = s.processInstallmentAndUpdateLimit(ctx, processInstallmentAndUpdateLimitInput{
				totalAmount:    validateOutput.totalAmount,
				limit:          limitOutput,
				contractNumber: contractNumber,
				startingTenor:  time.Now().UTC().AddDate(0, 1, 0).Month(),
				tx:             tx,
			})
			return
		})

		if err = erg.Wait(); err != nil {
			return tracer.Error(err)
		}
		return
	})
	return
}

func (s *service) validateCreateTransaction(input validateCreateInput) (output validateCreateOutput, err error) {
	stockMap := make(map[int64]products.GetAllOutputItem)
	for _, item := range input.products {
		stockMap[item.ID] = item
	}

	totalAmount := float64(0)

	for _, product := range input.createInput.Products {
		availableQty, exists := stockMap[product.ID]
		if !exists {
			return output, tracer.Error(ErrProductNotFound)
		}
		if availableQty.Qty < product.Qty {
			return output, tracer.Error(ErrStockProductNotAvailable)
		}

		totalAmount += availableQty.Price * float64(product.Qty)
	}

	if totalAmount > input.limit.RemainingAmount {
		return output, tracer.Error(ErrExceedLimit)
	}

	output = validateCreateOutput{
		totalAmount: totalAmount,
	}

	return
}

func (s *service) processInsertDataTransactionAndUpdateProduct(ctx context.Context,
	input processInsertDataTransactionAndUpdateProductInput,
) (transactionID int64, err error) {
	transactionCreateOutput, err := s.transactionRepository.Create(ctx, transactions.CreateInput{
		Transaction:     input.tx,
		LimitID:         input.limitID,
		ConsumerID:      input.consumerID,
		ContractNumber:  input.contractNumber,
		Amount:          input.totalAmount,
		TransactionDate: time.Now().UTC(),
		Status:          string(primitive.TransactionStatusActive),
	})
	if err != nil {
		return transactionID, tracer.Error(err)
	}
	transactionID = transactionCreateOutput.ID

	orderProducts := make(map[int64]int64)
	for _, product := range input.createInput.Products {
		orderProducts[product.ID] = product.Qty
	}

	transactionItems := make([]transaction_items.CreatesItemInput, 0)
	updateProducts := make([]products.UpdatesInputItem, 0)
	for _, product := range input.products {
		transactionItems = append(transactionItems, transaction_items.CreatesItemInput{
			MerchantID: product.MerchantID,
			Name:       product.Name,
			Image:      product.Image,
			Qty:        product.Qty,
			UnitPrice:  product.Price,
			Amount:     float64(product.Qty) * product.Price,
		})

		updateProducts = append(updateProducts, products.UpdatesInputItem{
			ID:  product.ID,
			Qty: product.Qty - orderProducts[product.ID],
		})
	}

	var erg errgroup.Group

	erg.Go(func() (err error) {
		err = s.transactionItemRepository.Creates(ctx, transaction_items.CreatesInput{
			Transaction:   input.tx,
			TransactionID: transactionCreateOutput.ID,
			Items:         transactionItems,
		})
		if err != nil {
			return tracer.Error(err)
		}

		return
	})

	erg.Go(func() (err error) {
		err = s.productRepository.Updates(ctx, products.UpdatesInput{
			Transaction: input.tx,
			Items:       updateProducts,
		})
		if err != nil {
			return tracer.Error(err)
		}

		return
	})

	if err = erg.Wait(); err != nil {
		return transactionID, tracer.Error(err)
	}

	return
}

func (s *service) processInstallmentAndUpdateLimit(ctx context.Context, input processInstallmentAndUpdateLimitInput) (
	err error) {
	calculateOutput := util.CalculateMonthlyInstallments(util.CalculateMonthlyInstallmentsInput{
		Principal:          input.totalAmount,
		AnnualInterestRate: 2.0,
		Tenor:              input.limit.Tenor,
	})

	timeNow := time.Now().UTC()
	startingInstallmentDate := time.Date(
		timeNow.Year(), input.startingTenor, 1, 0, 0, 0, 0, time.UTC,
	)

	installmentItems := []installments.CreatesInputItem{
		{
			Amount:  calculateOutput.MonthlyInstallments,
			DueDate: startingInstallmentDate,
			Status:  string(primitive.InstallmentStatusUnPaid),
		},
	}

	for i := int32(1); i < input.limit.Tenor; i++ {
		nextInstallmentDate := startingInstallmentDate.AddDate(0, int(i), 0)
		installmentItems = append(installmentItems, installments.CreatesInputItem{
			Amount:  calculateOutput.MonthlyInstallments,
			DueDate: nextInstallmentDate,
			Status:  string(primitive.InstallmentStatusUnPaid),
		})
	}

	var erg errgroup.Group

	erg.Go(func() (err error) {
		err = s.installmentRepository.Creates(ctx, installments.CreatesInput{
			Transaction:    input.tx,
			LimitID:        input.limit.ID,
			ContractNumber: input.contractNumber,
			Items:          installmentItems,
		})
		if err != nil {
			return tracer.Error(err)
		}

		return
	})

	erg.Go(func() (err error) {
		err = s.limitRepository.Update(ctx, limits.UpdateInput{
			Transaction:     input.tx,
			ID:              input.limit.ID,
			RemainingAmount: input.limit.RemainingAmount - input.totalAmount,
		})
		if err != nil {
			return tracer.Error(err)
		}

		return
	})

	if err = erg.Wait(); err != nil {
		return tracer.Error(err)
	}

	return
}
