package db

type Locking string

const (
	LockingUpdate Locking = "update"
	LockingDelete Locking = "delete"
)
