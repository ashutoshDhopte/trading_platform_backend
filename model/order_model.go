package model

type OrderModel struct {
	OrderID                int64
	StockTicker            string
	StockName              string
	TradeType              string
	OrderStatus            string
	Quantity               int64
	PricePerShareDollars   float64
	TotalOrderValueDollars float64
	CreatedAt              string
	Notes                  string
}
