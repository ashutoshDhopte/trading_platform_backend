package orm

import (
	"time"
)

type Stocks struct {
	StockID                int64 `gorm:"primaryKey"`
	Ticker                 string
	Name                   string
	OpeningPriceCents      float64
	MinPriceGeneratorCents float64
	MaxPriceGeneratorCents float64
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
