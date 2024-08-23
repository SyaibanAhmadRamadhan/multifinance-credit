package bank_account

import "github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"

type CreatesInput struct {
	ConsumerID int64
	Items      []CreatesInputItem
}

type CreatesInputItem struct {
	Name              string
	AccountNumber     string
	AccountHolderName string
}

type GetAllInput struct {
	ConsumerID int64
	Pagination pagination.PaginationInput
}

type GetAllOutput struct {
	Pagination pagination.PaginationOutput
	Items      []GetAllOutputItem
}

type GetAllOutputItem struct {
	ID                int64
	ConsumerID        int64
	Name              string
	AccountNumber     string
	AccountHolderName string
}
