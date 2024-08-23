package primitive

type InstallmentStatus string

const (
	InstallmentStatusUnPaid    InstallmentStatus = "UNPAID"
	InstallmentStatusPaid      InstallmentStatus = "PAID"
	InstallmentStatusDefaulted InstallmentStatus = "DEFAULTED"
)
