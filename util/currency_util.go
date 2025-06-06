package util

import (
	"fmt"
	"strconv"
)

const InitialInvestmentCents = 10000000

func ConvertCentsToDollars(cents int64) float64 {

	dollars := cents / 100.0
	rounded := fmt.Sprintf("%d", dollars)
	parsedDollar, err := strconv.ParseFloat(rounded, 64)
	if err != nil {
		return 0
	}
	return parsedDollar
}
