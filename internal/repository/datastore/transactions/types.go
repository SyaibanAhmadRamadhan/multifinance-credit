package transactions

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"time"
)

type CreateInput struct {
	Transaction     db.Rdbms
	LimitID         int64
	ConsumerID      int64
	ContractNumber  int64
	Amount          float64
	TransactionDate time.Time
	Status          string
}

type CreateOutput struct {
	ID int64
}
