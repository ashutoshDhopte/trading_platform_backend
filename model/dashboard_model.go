package model

type DashboardModel struct {
	User                     UserModel
	Stocks                   []StockModel
	Holdings                 []HoldingModel
	TotalHoldingValueDollars float64
	PortfolioValueDollars    float64
	TotalPnLDollars          float64
	TotalReturnPercent       float64
}
