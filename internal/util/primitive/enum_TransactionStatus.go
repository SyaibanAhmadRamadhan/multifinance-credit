package primitive

type TransactionStatus string

const (
	TransactionStatusActive    TransactionStatus = "ACTIVE"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
)
