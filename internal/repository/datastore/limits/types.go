package limits

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/guregu/null/v5"
)

type GetInput struct {
	Tx         db.Rdbms
	Locking    db.Locking
	ID         null.Int
	Tenor      null.Int32
	ConsumerID null.Int
}

type GetOutput struct {
	ID              int64   `db:"id"`
	ConsumerID      int64   `db:"consumer_id"`
	Tenor           int32   `db:"tenor"`
	Amount          float64 `db:"amount"`
	RemainingAmount float64 `db:"remaining_amount"`
}

type CreatesInput struct {
	Transaction db.Rdbms
	ConsumerID  int64
	Items       []CreatesInputItem
}

type CreatesInputItem struct {
	Tenor  int32
	Amount float64
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
	ID         int64   `db:"id"`
	ConsumerID int64   `db:"consumer_id"`
	Tenor      int32   `db:"tenor"`
	Amount     float64 `db:"amount"`
}

type UpdateInput struct {
	Transaction     db.Rdbms
	ID              int64
	RemainingAmount float64
}
