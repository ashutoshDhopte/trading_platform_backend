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

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func NewsMigration() {

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		stocks := db.GetAllStocks()

		for i, _ := range stocks {
			newsArr := finnhub2.FetchNewsFromFinnhubMigration(stocks[i].Ticker)

			if len(newsArr) == 0 {
				continue
			}

			newsArticleMap := make(map[string]*orm.NewsArticles)

			for _, news := range newsArr {

				newsArticle := orm.NewsArticles{
					Ticker:          stocks[i].Ticker,
					ArticleTitle:    news.GetHeadline(),
					ArticleSummary:  news.GetSummary(),
					ArticleURL:      news.GetUrl(),
					PublicationTime: time.Unix(news.GetDatetime(), 0),
				}

				newsArticleMap[newsArticle.ArticleTitle] = &newsArticle
			}

			sentimentResponses, err := getArticlesSentiment(newsArr)
			if err != nil {
				panic(err)
			}

			for _, sentiment := range sentimentResponses {
				newsArticleMap[sentiment.ArticleTitle].SentimentScore = sentiment.Score
				err := db.DB.Create(newsArticleMap[sentiment.ArticleTitle]).Error
				if err != nil {
					panic(err)
				}
			}

			// Calculate 14-day EMA for sentiment
			emaScore := calculateSentimentEMA(stocks[i].StockID, 14)
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

// calculateSentimentEMA calculates the Exponential Moving Average for sentiment scores
func calculateSentimentEMA(stockId int64, days int) float32 {
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

// UpdateSentimentEMA updates the sentiment EMA for all stocks
func UpdateSentimentEMA() {
	stocks := db.GetAllStocks()

	for _, stock := range stocks {
		emaScore := calculateSentimentEMA(stock.StockID, 14)

		err := db.DB.Model(&stock).Select("overall_sentiment_score").Updates(map[string]interface{}{
			"overall_sentiment_score": emaScore,
		}).Error

		if err != nil {
			fmt.Printf("Error updating sentiment EMA for stock %s: %v\n", stock.Ticker, err)
		} else {
			fmt.Printf("Updated sentiment EMA for %s: %.4f\n", stock.Ticker, emaScore)
		}
	}
}

func getArticlesSentiment(companyNewsList []finnhub.CompanyNews) ([]model.SentimentResponse, error) {

	_ = godotenv.Load()
	sentimentAnalysisModelUrl := os.Getenv("SENTIMENT_ANALYSIS_MODEL_URL")
	url := sentimentAnalysisModelUrl + "/sentiment/analyze"
	fmt.Println("Calling API endpoint:", url)

	sentimentReqList := make([]model.SentimentRequest, 0)
	var sentimentResponses []model.SentimentResponse

	for _, news := range companyNewsList {

		sentimentReqList = append(sentimentReqList, model.SentimentRequest{
			ArticleTitle:   news.GetHeadline(),
			ArticleSummary: news.GetSummary(),
		})
	}

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
