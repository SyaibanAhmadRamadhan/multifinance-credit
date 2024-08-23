package products

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/guregu/null/v5"
)

type GetInput struct {
	ID null.Int
}

type GetOutput struct {
	ID         int64   `db:"id"`
	MerchantID int64   `db:"merchant_id"`
	Image      string  `db:"image"`
	Name       string  `db:"name"`
	Qty        int64   `db:"qty"`
	Price      float64 `db:"price"`
}

type CreatesInput struct {
	Transaction *db.SqlxWrapper
	MerchantID  int64
	Items       []CreatesInputItem
}

type CreatesInputItem struct {
	Name  string
	Image string
	Qty   int64
	Price float64
}

type GetAllInput struct {
	Locking     db.Locking
	Transaction *db.SqlxWrapper
	MerchantID  null.Int
	IDs         []int64
	Pagination  pagination.PaginationInput
}

type GetAllOutput struct {
	Pagination pagination.PaginationOutput
	Items      []GetAllOutputItem
}

type GetAllOutputItem struct {
	ID         int64   `db:"id"`
	MerchantID int64   `db:"merchant_id"`
	Image      string  `db:"image"`
	Name       string  `db:"name"`
	Qty        int64   `db:"qty"`
	Price      float64 `db:"price"`
}

type UpdatesInput struct {
	Transaction *db.SqlxWrapper
	Items       []UpdatesInputItem
}

type UpdatesInputItem struct {
	ID  int64
	Qty int64
}
