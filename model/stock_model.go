package model

type StockModel struct {
	StockID             int64
	Ticker              string
	Name                string
	OpeningPriceDollars float64
	ChangedPriceDollars float64
	ChangedPercent      float64
	UpdatedAt           string
}
