package orm

import (
	"time"
)

type Orders struct {
	OrderID              int64 `gorm:"primaryKey"`
	UserID               int64
	StockID              int64
	TradeType            string
	OrderStatus          string
	Quantity             int64
	PricePerShareCents   int64
	TotalOrderValueCents int64
	CreatedAt            time.Time
	Notes                string
}
