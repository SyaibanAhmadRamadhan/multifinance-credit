package util

import "math"

type CalculateMonthlyInstallmentsInput struct {
	Principal          float64 // Principal is loan amount
	AnnualInterestRate float64
	Tenor              int64
}

type CalculateMonthlyInstallmentsOutput struct {
	MonthlyInstallments float64
}

func CalculateMonthlyInstallments(input CalculateMonthlyInstallmentsInput) CalculateMonthlyInstallmentsOutput {
	monthlyInterestRate := (input.AnnualInterestRate / 12) / 100

	monthlyInstallment := input.Principal * monthlyInterestRate * math.Pow(1+monthlyInterestRate, float64(input.Tenor)) /
		(math.Pow(1+monthlyInterestRate, float64(input.Tenor)) - 1)

	return CalculateMonthlyInstallmentsOutput{
		MonthlyInstallments: monthlyInstallment,
	}
}
