package orm

import (
	"time"
)

type Holdings struct {
	HoldingID                int64 `gorm:"primaryKey"`
	UserID                   int64
	StockID                  int64
	Quantity                 int64
	AverageCostPerShareCents int64
	CreatedAt                time.Time
	UpdatedAt                time.Time
}
