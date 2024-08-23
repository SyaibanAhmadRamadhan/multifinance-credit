package product

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/primitive"
	"github.com/guregu/null/v5"
)

type CreateInput struct {
	MerchantID int64
	Name       string
	Image      primitive.PresignedFileUpload
	Qty        int64
	Price      float64
}

type CreateOutput struct {
	Image primitive.PresignedFileUploadOutput
}

type GetAllInput struct {
	MerchantID null.Int
	IDs        []int64
	Pagination pagination.PaginationInput
}

type GetAllOutput struct {
	Pagination pagination.PaginationOutput
	Items      []GetAllOutputItem
}

type GetAllOutputItem struct {
	ID         int64
	MerchantID int64
	Name       string
	Image      string
	Qty        int64
	Price      float64
}
