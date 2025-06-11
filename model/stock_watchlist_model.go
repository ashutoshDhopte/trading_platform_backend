package model

import "time"

type StockWatchlistModel struct {
	StockWatchlistID int64
	UserId           int64
	StockId          int64
	TargetPriceCents int64
	IsActive         bool
	CreatedAt        time.Time
}
