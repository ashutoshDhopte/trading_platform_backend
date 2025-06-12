package model

type StockWatchlistModel struct {
	StockWatchlistID   int32
	UserId             int32
	StockId            int32
	StockTicker        string
	StockName          string
	TargetPriceDollars float64
	DiffPriceDollars   float64
	DiffPercent        float64
	IsActive           bool
	CreatedAt          string
}
