package transaction_items

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
)

type CreatesInput struct {
	Transaction   db.Rdbms
	TransactionID int64
	Items         []CreatesItemInput
}
type CreatesItemInput struct {
	MerchantID int64
	Name       string
	Image      string
	Qty        int64
	UnitPrice  float64
	Amount     float64
}
