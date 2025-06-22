package orm

import "time"

type NewsArticles struct {
	NewsArticleID   int `gorm:"primary_key"`
	Ticker          string
	ArticleTitle    string
	ArticleURL      string
	PublicationTime time.Time
	SentimentScore  float32
}
