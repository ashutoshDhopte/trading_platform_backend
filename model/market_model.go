package model

type MarketModel struct {
	Stock StockModel
	News  []NewsModel
}

type NewsModel struct {
	NewsArticleID   int     `json:"newsArticleId"`
	Ticker          string  `json:"ticker"`
	ArticleTitle    string  `json:"articleTitle"`
	ArticleSummary  string  `json:"articleSummary"`
	ArticleURL      string  `json:"articleUrl"`
	PublicationTime string  `json:"publicationTime"`
	SentimentScore  float32 `json:"sentimentScore"`
}
