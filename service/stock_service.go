package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"trading_platform_backend/db"
	finnhub2 "trading_platform_backend/external_client"
	"trading_platform_backend/model"
	"trading_platform_backend/orm"
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
					ArticleURL:      news.GetUrl(),
					PublicationTime: time.Unix(news.GetDatetime(), 0),
				}

				newsArticleMap[newsArticle.ArticleTitle] = &newsArticle
			}

			sentimentResponses, err := getArticlesSentiment(newsArr)
			if err != nil {
				panic(err)
			}

			totalSentimentScore := float32(0)

			for _, sentiment := range sentimentResponses {
				newsArticleMap[sentiment.ArticleTitle].SentimentScore = sentiment.Score
				//err := db.DB.Create(newsArticleMap[sentiment.ArticleTitle]).Error
				//if err != nil {
				//	panic(err)
				//}
				totalSentimentScore += sentiment.Score
			}

			stocks[i].OverallSentimentScore = totalSentimentScore / float32(len(sentimentResponses))

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
