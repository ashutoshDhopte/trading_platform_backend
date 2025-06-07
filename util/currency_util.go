package util

const InitialInvestmentCents = 10000000

func ConvertCentsToDollars(cents int64) float64 {
	dollars := float64(cents) / 100.0
	truncated := float64(int(dollars*100)) / 100.0
	return truncated
}
