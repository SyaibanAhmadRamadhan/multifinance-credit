package bank_account

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/bank_accounts"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/consumers"
)

type service struct {
	bankAccountRepository bank_accounts.Repository
	consumerRepository    consumers.Repository
	dbTx                  db.SqlxTransaction
}

var _ Service = (*service)(nil)

type NewServiceOpts struct {
	BankAccountRepository bank_accounts.Repository
	ConsumerRepository    consumers.Repository
	DBTx                  db.SqlxTransaction
}

func NewService(
	opts NewServiceOpts,
) *service {
	return &service{
		consumerRepository:    opts.ConsumerRepository,
		bankAccountRepository: opts.BankAccountRepository,
		dbTx:                  opts.DBTx,
	}
}
