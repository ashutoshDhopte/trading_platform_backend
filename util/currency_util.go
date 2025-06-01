package util

import (
	"fmt"
	"strconv"
)

func ConvertCentsToDollars(cents float64) float64 {

	dollars := cents / 100.0
	rounded := fmt.Sprintf("%.2f", dollars)
	parsedDollar, err := strconv.ParseFloat(rounded, 64)
	if err != nil {
		return 0
	}
	return parsedDollar
}
