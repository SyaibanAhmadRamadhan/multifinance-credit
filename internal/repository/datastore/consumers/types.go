package consumers

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/db"
	"github.com/guregu/null/v5"
	"time"
)

type CreateInput struct {
	Transaction  *db.SqlxWrapper
	UserID       int64
	Nik          string
	FullName     string
	LegalName    string
	PlaceOfBirth string
	DateOfBirth  time.Time
	Salary       float64
	PhotoKTP     string
	PhotoSelfie  string
}

type CreateOutput struct {
	ID int64
}

type CheckExistingInput struct {
	ByID  null.Int64
	ByNIK null.String
}

type CheckExistingOutput struct {
	Existing bool
}

type GetInput struct {
	ID     null.Int
	UserID null.Int
}

type GetOutput struct {
	ID           int64     `db:"id"`
	UserID       int64     `db:"user_id"`
	Nik          string    `db:"nik"`
	FullName     string    `db:"full_name"`
	LegalName    string    `db:"legal_name"`
	PlaceOfBirth string    `db:"place_of_birth"`
	DateOfBirth  time.Time `db:"date_of_birth"`
	Salary       float64   `db:"salary"`
	PhotoKTP     string    `db:"photo_ktp"`
	PhotoSelfie  string    `db:"photo_selfie"`
}
