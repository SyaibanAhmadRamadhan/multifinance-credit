package util_test

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util"
	"math"
	"testing"
)

func TestCalculateMonthlyInstallments(t *testing.T) {
	tests := []struct {
		name                string
		principal           float64
		annualInterestRate  float64
		tenorInMonths       int32
		expectedInstallment float64
	}{
		{
			name:                "Loan 10 million, interest 12%, tenor 24 months",
			principal:           10000000.0,
			annualInterestRate:  12.0,
			tenorInMonths:       24,
			expectedInstallment: 470734.72,
		},
		{
			name:                "5 million loan, 10% interest, 12 month term",
			principal:           5000000.0,
			annualInterestRate:  10.0,
			tenorInMonths:       12,
			expectedInstallment: 439579.44,
		},
		{
			name:                "20 million loan, 15% interest, 36 month term",
			principal:           20000000.0,
			annualInterestRate:  15.0,
			tenorInMonths:       36,
			expectedInstallment: 693306.57,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualInstallment := util.CalculateMonthlyInstallments(util.CalculateMonthlyInstallmentsInput{
				Principal:          tt.principal,
				AnnualInterestRate: tt.annualInterestRate,
				Tenor:              tt.tenorInMonths,
			})
			if math.Abs(actualInstallment.MonthlyInstallments-tt.expectedInstallment) > 0.01 {
				t.Errorf("got %.2f, want %.2f", actualInstallment, tt.expectedInstallment)
			}
		})
	}
}
