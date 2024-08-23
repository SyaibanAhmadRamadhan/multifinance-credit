package util

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type CalculateMonthlyInstallmentsInput struct {
	Principal          float64 // Principal is loan amount
	AnnualInterestRate float64
	Tenor              int32
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

func GenerateContractNumber() int64 {
	randSource := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	randomNumber := randSource.Intn(100000)

	timestamp := time.Now().Format("20060102150405")

	contractNumber := fmt.Sprintf("%s%04d", timestamp, randomNumber)

	contractNumberInt, err := strconv.ParseInt(contractNumber, 10, 64)
	Panic(err)
	return contractNumberInt
}
