package bank_accounts

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/guregu/null/v5"
)

type GetInput struct {
	ID            null.Int
	AccountNumber null.String
}

type GetOutput struct {
	ID                int64  `db:"id"`
	ConsumerID        int64  `db:"consumer_id"`
	Name              string `db:"name"`
	AccountNumber     string `db:"account_number"`
	AccountHolderName string `db:"account_holder_name"`
}

type CreatesInput struct {
	Transaction *db.SqlxWrapper
	Items       []CreatesInputItem
}

type CreatesInputItem struct {
	ConsumerID        int64
	Name              string
	AccountNumber     string
	AccountHolderName string
}

type GetAllInput struct {
	ConsumerID null.Int
	Pagination pagination.PaginationInput
}

type GetAllOutput struct {
	Pagination pagination.PaginationOutput
	Items      []GetAllOutputItem
}

type GetAllOutputItem struct {
	ID                int64  `db:"id"`
	ConsumerID        int64  `db:"consumer_id"`
	Name              string `db:"name"`
	AccountNumber     string `db:"account_number"`
	AccountHolderName string `db:"account_holder_name"`
}
