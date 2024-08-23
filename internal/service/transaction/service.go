package transaction

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/installments"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/limits"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/products"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/transaction_items"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/transactions"
)

type service struct {
	productRepository         products.Repository
	transactionRepository     transactions.Repository
	transactionItemRepository transaction_items.Repository
	installmentRepository     installments.Repository
	limitRepository           limits.Repository
	dbTx                      db.SqlxTransaction
}

var _ Service = (*service)(nil)

type NewServiceOpts struct {
	ProductRepository         products.Repository
	TransactionRepository     transactions.Repository
	TransactionItemRepository transaction_items.Repository
	InstallmentRepository     installments.Repository
	LimitRepository           limits.Repository
	DBTx                      db.SqlxTransaction
}

func NewService(
	opts NewServiceOpts,
) *service {
	return &service{
		productRepository:         opts.ProductRepository,
		transactionRepository:     opts.TransactionRepository,
		transactionItemRepository: opts.TransactionItemRepository,
		installmentRepository:     opts.InstallmentRepository,
		limitRepository:           opts.LimitRepository,
		dbTx:                      opts.DBTx,
	}
}
