package orm

import "time"

type StockWatchlist struct {
	StockWatchlistID int64 `gorm:"primaryKey"`
	UserId           int64
	StockId          int64
	TargetPriceCents int64
	IsActive         bool
	CreatedAt        time.Time
}

func (StockWatchlist) TableName() string {
	return "stock_watchlist" // replace with your desired table name
}
