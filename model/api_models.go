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
