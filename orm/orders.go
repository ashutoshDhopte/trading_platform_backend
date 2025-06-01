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
	PricePerShareCents   float64
	TotalOrderValueCents float64
	CreatedAt            time.Time
	Notes                string
}
