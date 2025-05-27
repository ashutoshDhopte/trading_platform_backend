package model

type DashboardModel struct {
	User     UserModel
	Stocks   []StockModel
	Holdings []HoldingModel
}
