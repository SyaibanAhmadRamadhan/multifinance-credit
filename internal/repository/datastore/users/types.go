package users

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/guregu/null/v5"
)

type CheckExistingInput struct {
	ByID    null.Int
	ByEmail null.String
}

type CheckExistingOutput struct {
	Existing bool
}

type CreateInput struct {
	Transaction db.Rdbms
	Email       string
	Password    string
}

type CreateOutput struct {
	ID int64
}

type GetInput struct {
	ID    null.Int
	Email null.String
}

type GetOutput struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
}
