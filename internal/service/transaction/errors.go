package transaction

import (
	"errors"
)

var ErrStockProductNotAvailable = errors.New("stock product not available")

var ErrProductNotFound = errors.New("product not found")

var ErrLimitNotFound = errors.New("limit not found")
var ErrExceedLimit = errors.New("exceeds the limit")
