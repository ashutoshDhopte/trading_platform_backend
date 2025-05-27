package util

import (
	"fmt"
	"strconv"
)

func ConvertCentsToDollars(cents int64) float64 {

	dollars := float64(cents) / 100.0
	rounded := fmt.Sprintf("%.2f", dollars)
	parsedDollar, err := strconv.ParseFloat(rounded, 64)
	if err != nil {
		return 0
	}
	return parsedDollar
}
