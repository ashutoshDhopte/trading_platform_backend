package model

type ApiResponse struct {
	Success      bool
	Data         interface{}
	ErrorMessage string
}

type SentimentRequest struct {
	FinnhubNewsID  string `json:"finnhub_news_id"`
	ArticleTitle   string `json:"articleTitle"`
	ArticleSummary string `json:"articleSummary"`
}

type SentimentResponse struct {
	FinnhubNewsID string  `json:"finnhub_news_id"`
	Score         float32 `json:"score"`
}

type After struct {
	UserID             int    `json:"user_id"`
	StockID            int    `json:"stock_id"`
	TradeType          string `json:"trade_type"`
	Quantity           int    `json:"quantity"`
	PricePerShareCents int    `json:"price_per_share_cents"`
}

type SocialChanMessage struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}

type WebSocketMessage struct {
	EventType string      `json:"event_type"`
	Content   interface{} `json:"content"`
}
