package model

type HoldingModel struct {
	HoldingID                  int64
	StockTicker                string
	Quantity                   int64
	AverageCostPerShareDollars float64
	TotalValueDollars          float64
	PnLDollars                 float64
	PnLPercent                 float64
	UpdatedAt                  string
}
