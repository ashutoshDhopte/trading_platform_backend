package orm

import "time"

type StockWatchlist struct {
	StockWatchlistID int32 `gorm:"primaryKey"`
	UserId           int32
	StockId          int32
	TargetPriceCents int64
	IsActive         bool
	CreatedAt        time.Time
}

func (StockWatchlist) TableName() string {
	return "stock_watchlist" // replace with your desired table name
}
