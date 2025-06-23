package model

type StockModel struct {
	StockID               int64
	Ticker                string
	Name                  string
	OpeningPriceDollars   float64
	CurrentPriceDollars   float64
	ChangedPriceDollars   float64
	ChangedPercent        float64
	UpdatedAt             string
	OverallSentimentScore float32
}

func (stock *StockModel) GetChangedPriceDollars() float64 {
	return stock.CurrentPriceDollars - stock.OpeningPriceDollars
}

func (stock *StockModel) GetChangedPercent() float64 {
	return stock.GetChangedPriceDollars() / stock.OpeningPriceDollars
}
