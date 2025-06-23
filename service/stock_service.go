package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"trading_platform_backend/db"
	finnhub2 "trading_platform_backend/external_client"
	"trading_platform_backend/model"
	"trading_platform_backend/orm"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func NewsMigration() {

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		stocks := db.GetAllStocks()

		for i, _ := range stocks {
			// For migration, fetch a wide date range
			newsArr := finnhub2.FetchNewsFromFinnhub(stocks[i].Ticker, "2025-03-01", "2025-06-21")

			if len(newsArr) == 0 {
				continue
			}

			newsArticleMap := make(map[string]*orm.NewsArticles)
			var sentimentReqs []model.SentimentRequest

			for _, news := range newsArr {
				newsID := strconv.FormatInt(news.GetId(), 10)
				var count int64
				db.DB.Model(&orm.NewsArticles{}).Where("finnhub_news_id = ?", newsID).Count(&count)
				if count > 0 {
					continue
				}

				newsArticle := orm.NewsArticles{
					Ticker:          stocks[i].Ticker,
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

			sentimentResponses, err := GetArticlesSentiment(sentimentReqs)
			if err == nil {
				for _, sentiment := range sentimentResponses {
					if article, ok := newsArticleMap[sentiment.FinnhubNewsID]; ok {
						article.SentimentScore = sentiment.Score
					}
				}
			}

			for _, article := range newsArticleMap {
				err := db.DB.Create(article).Error
				if err != nil {
					panic(err)
				}
			}

			// Calculate 14-day EMA for sentiment
			emaScore := CalculateSentimentEMA(stocks[i].StockID, 14)
			stocks[i].OverallSentimentScore = emaScore

			err = tx.Model(&stocks[i]).Select("overall_sentiment_score").Updates(&stocks[i]).Error
			if err != nil {
				panic(err)
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

// CalculateSentimentEMA calculates the Exponential Moving Average for sentiment scores
func CalculateSentimentEMA(stockId int64, days int) float32 {
	// Get sentiment data for the specified number of days
	sentimentData := db.GetSentimentDataForStock(stockId, days)

	if len(sentimentData) == 0 {
		return 0.0
	}

	// Calculate the multiplier for EMA (2 / (n + 1) where n = days)
	multiplier := 2.0 / float64(days+1)

	// Start with the first sentiment score as the initial EMA
	ema := float64(sentimentData[0].SentimentScore)

	// Calculate EMA for the remaining data points
	for i := 1; i < len(sentimentData); i++ {
		currentSentiment := float64(sentimentData[i].SentimentScore)
		ema = (currentSentiment * multiplier) + (ema * (1 - multiplier))
	}

	return float32(ema)
}

// GetArticlesSentiment calls the sentiment analysis API for a batch of articles
func GetArticlesSentiment(sentimentReqList []model.SentimentRequest) ([]model.SentimentResponse, error) {
	_ = godotenv.Load()
	sentimentAnalysisModelUrl := os.Getenv("SENTIMENT_ANALYSIS_MODEL_URL")
	url := sentimentAnalysisModelUrl + "/sentiment/analyze"
	fmt.Println("Calling API endpoint:", url)

	var sentimentResponses []model.SentimentResponse

	jsonData, err1 := json.Marshal(sentimentReqList)
	if err1 != nil {
		fmt.Printf("Error marshaling to JSON: %s\n", err1)
		return sentimentResponses, err1
	}

	resp, err2 := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err2 != nil {
		fmt.Printf("Error sending request: %s\n", err2)
		return sentimentResponses, err2
	}

	defer resp.Body.Close()

	fmt.Printf("\nReceived response with status code: %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API call failed with status code: %d\n", resp.StatusCode)
		return sentimentResponses, errors.New("API call failed with status code: " + strconv.Itoa(resp.StatusCode))
	}

	responseBody, err3 := io.ReadAll(resp.Body)
	if err3 != nil {
		fmt.Printf("Error reading response body: %s\n", err3)
		return sentimentResponses, err3
	}

	err3 = json.Unmarshal(responseBody, &sentimentResponses)
	if err3 != nil {
		fmt.Printf("Error unmarshaling response JSON: %s\n", err3)
		return sentimentResponses, err3
	}

	return sentimentResponses, nil
}
