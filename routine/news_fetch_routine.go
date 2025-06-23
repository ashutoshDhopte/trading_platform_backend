package routine

import (
	"fmt"
	"strconv"
	"time"
	"trading_platform_backend/db"
	"trading_platform_backend/external_client"
	"trading_platform_backend/model"
	"trading_platform_backend/orm"
	"trading_platform_backend/service"
)

func initNewsFetchRoutine() {
	go startNewsFetchLoop()
}

func startNewsFetchLoop() {
	// Run every 15 minutes
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for {
		fetchAndSaveLatestNewsForAllStocks()
		<-ticker.C
	}
}

func fetchAndSaveLatestNewsForAllStocks() {
	stocks := db.GetAllStocks()
	to := time.Now().UTC()
	from := to.Add(-16 * time.Minute) // fetch news from the last interval (with 1 min overlap)
	fromStr := from.Format("2006-01-02")
	toStr := to.Format("2006-01-02")

	for _, stock := range stocks {
		newsArr := external_client.FetchNewsFromFinnhub(stock.Ticker, fromStr, toStr)
		if len(newsArr) == 0 {
			continue
		}

		newsArticleMap := make(map[string]*orm.NewsArticles)
		var sentimentReqs []model.SentimentRequest
		for _, news := range newsArr {
			var count int64
			newsID := strconv.FormatInt(news.GetId(), 10)
			db.DB.Model(&orm.NewsArticles{}).Where("finnhub_news_id = ?", newsID).Count(&count)
			if count > 0 {
				continue
			}

			newsArticle := orm.NewsArticles{
				Ticker:          stock.Ticker,
				FinnhubNewsID:   newsID,
				ArticleTitle:    news.GetHeadline(),
				ArticleSummary:  news.GetSummary(),
				ArticleURL:      news.GetUrl(),
				PublicationTime: time.Unix(news.GetDatetime(), 0),
			}
			newsArticleMap[newsID] = &newsArticle
			sentimentReqs = append(sentimentReqs, model.SentimentRequest{
				FinnhubNewsID:  newsID,
				ArticleTitle:   news.GetHeadline(),
				ArticleSummary: news.GetSummary(),
			})
		}

		if len(newsArticleMap) == 0 {
			continue
		}

		sentimentResponses, err := service.GetArticlesSentiment(sentimentReqs)
		if err == nil {
			for _, sentiment := range sentimentResponses {
				if article, ok := newsArticleMap[sentiment.FinnhubNewsID]; ok {
					article.SentimentScore = sentiment.Score
				}
			}
		}

		for _, article := range newsArticleMap {
			db.DB.Create(article)
		}

		// Recalculate EMA for this stock
		emaScore := service.CalculateSentimentEMA(stock.StockID, 14)
		db.DB.Model(&stock).Select("overall_sentiment_score").Updates(map[string]interface{}{
			"overall_sentiment_score": emaScore,
		})
		fmt.Printf("[NewsRoutine] Updated EMA for %s: %.4f\n", stock.Ticker, emaScore)
	}
}
