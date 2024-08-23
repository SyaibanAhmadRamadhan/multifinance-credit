package transaction

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/limits"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/products"
	"time"
)

type CreateInput struct {
	Products   []CreateInputProductItem
	LimitID    int64
	ConsumerID int64
}

type CreateInputProductItem struct {
	ID  int64
	Qty int64
}

type CreateOutput struct {
	ID int64
}

type validateCreateInput struct {
	limit       limits.GetOutput
	products    []products.GetAllOutputItem
	createInput CreateInput
}

type validateCreateOutput struct {
	totalAmount float64
}

type processInsertDataTransactionAndUpdateProductInput struct {
	tx             *db.SqlxWrapper
	products       []products.GetAllOutputItem
	createInput    CreateInput
	contractNumber int64
	limitID        int64
	totalAmount    float64
	consumerID     int64
}

type processInstallmentAndUpdateLimitInput struct {
	totalAmount    float64
	limit          limits.GetOutput
	contractNumber int64
	startingTenor  time.Month
	tx             *db.SqlxWrapper
}
