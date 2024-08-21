package entity

import "time"

type User struct {
	ID        int64     `db:"id"`
	Password  string    `db:"password"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}
