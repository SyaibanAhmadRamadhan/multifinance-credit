package entity

import "time"

type Consumer struct {
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
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
