package orm

import (
	"time"
)

type StockOHLCV struct {
	StockOHLCVId int32 `gorm:"primary_key"`
	StockName    string
	Timestamp    time.Time
	Open         float32
	High         float32
	Low          float32
	Close        float32
	Volume       int32
}

func (StockOHLCV) TableName() string {
	return "stock_ohlcv" // replace with your desired table name
}
