package installments

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/guregu/null/v5"
	"time"
)

type GetInput struct {
	ID null.Int
}

type GetOutput struct {
	ID              int64      `db:"id"`
	LimitID         int64      `db:"limit_id"`
	PaymentMethodID *int64     `db:"payment_method_id"`
	ContractNumber  int64      `db:"contract_number"`
	Amount          float64    `db:"amount"`
	DueDate         time.Time  `db:"due_date"`
	PaymentDate     *time.Time `db:"payment_date"`
	Status          string     `db:"status"`
}

type CreatesInput struct {
	Transaction      *db.SqlxWrapper
	LimitID        int64
	ContractNumber int64
	Items          []CreatesInputItem
}

type CreatesInputItem struct {
	Amount  float64
	DueDate time.Time
	Status  string
}

type GetAllInput struct {
	ContractNumber int64
	Pagination     pagination.PaginationInput
}

type GetAllOutput struct {
	Pagination pagination.PaginationOutput
	Items      []GetAllOutputItem
}

type GetAllOutputItem struct {
	ID              int64      `db:"id"`
	LimitID         int64      `db:"limit_id"`
	PaymentMethodID *int64     `db:"payment_method_id"`
	ContractNumber  int64      `db:"contract_number"`
	Amount          float64    `db:"amount"`
	DueDate         time.Time  `db:"due_date"`
	PaymentDate     *time.Time `db:"payment_date"`
	Status          string     `db:"status"`
}
