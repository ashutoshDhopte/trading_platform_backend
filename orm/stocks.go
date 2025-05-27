package orm

import (
	"time"
)

type Stocks struct {
	StockID                int64 `gorm:"primaryKey"`
	Ticker                 string
	Name                   string
	OpeningPriceCents      int64
	MinPriceGeneratorCents int64
	MaxPriceGeneratorCents int64
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
