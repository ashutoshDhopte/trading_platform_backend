package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func MigrateStockOHLCV(ticker string) error {

	// CandleData represents the open, high, low, close, and volume for a specific timestamp.
	type CandleData struct {
		Open   string `json:"1. open"`
		High   string `json:"2. high"`
		Low    string `json:"3. low"`
		Close  string `json:"4. close"`
		Volume string `json:"5. volume"`
	}

	// TimeSeriesData represents the nested structure for a specific time series (e.g., 1-minute data).
	// The key (e.g., "2020-06-30 19:44:00") is a string, and the value is CandleData.
	type TimeSeriesData map[string]CandleData

	// AlphaVantageResponse represents the top-level structure of your JSON data.
	type AlphaVantageResponse struct {
		MetaData       map[string]interface{} `json:"Meta Data"`
		TimeSeries1min TimeSeriesData         `json:"Time Series (1min)"`
	}

	_ = godotenv.Load()
	alphaVintageToken := os.Getenv("ALPHA_VINTAGE_TOKEN")
	alphaVantageUrl := "https://alphavantage.co/query"

	apiCall := func(fullURL string) (AlphaVantageResponse, error) {

		fmt.Println("Calling API endpoint:", fullURL)

		var data AlphaVantageResponse

		resp, err2 := http.Get(fullURL)
		if err2 != nil {
			fmt.Printf("Error sending request: %s\n", err2)
			return data, err2
		}

		defer resp.Body.Close()

		fmt.Printf("Received response with status code: %d\n", resp.StatusCode)
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("API call failed with status code: %d\n", resp.StatusCode)
			return data, errors.New("API call failed with status code: " + strconv.Itoa(resp.StatusCode))
		}

		responseBody, err3 := io.ReadAll(resp.Body)
		if err3 != nil {
			fmt.Printf("Error reading response body: %s\n", err3)
			return data, err3
		}

		err3 = json.Unmarshal(responseBody, &data)
		if err3 != nil {
			fmt.Printf("Error unmarshaling response JSON: %s\n", err3)
			return data, err3
		}

		if len(data.TimeSeries1min) == 0 {
			fmt.Println(string(responseBody))
			return data, errors.New(string(responseBody))
		}

		return data, nil
	}

	currentTime := time.Now()

	for i := 12; i <= 60; i++ {

		err := db.DB.Transaction(func(tx *gorm.DB) error {

			timeStr := currentTime.AddDate(0, -i, 0).Format("2006-01")

			fmt.Println(timeStr)

			params := url.Values{}
			params.Add("function", "TIME_SERIES_INTRADAY")
			params.Add("symbol", ticker)
			params.Add("interval", "1min")
			params.Add("outputsize", "full")
			params.Add("apikey", alphaVintageToken)
			params.Add("month", timeStr)
			// Encode the parameters into a query string
			queryString := params.Encode()

			// Construct the full URL
			fullURL := fmt.Sprintf("%s?%s", alphaVantageUrl, queryString)

			data, err := apiCall(fullURL)
			if err != nil {
				return err
			}

			if len(data.TimeSeries1min) > 0 {
				for timestamp, candle := range data.TimeSeries1min {
					fmt.Println(timestamp)

					stockOHLCV := orm.StockOHLCV{StockName: ticker}
					stockOHLCV.Timestamp, err = time.Parse("2006-01-02 15:04:05", timestamp)
					parsed, err := strconv.ParseFloat(candle.Open, 32)
					if err == nil {
						stockOHLCV.Open = float32(parsed)
					} else {
						return err
					}

					parsed, err = strconv.ParseFloat(candle.Open, 32)
					if err == nil {
						stockOHLCV.Open = float32(parsed)
					} else {
						return err
					}

					parsed, err = strconv.ParseFloat(candle.High, 32)
					if err == nil {
						stockOHLCV.High = float32(parsed)
					} else {
						return err
					}

					parsed, err = strconv.ParseFloat(candle.Low, 32)
					if err == nil {
						stockOHLCV.Low = float32(parsed)
					} else {
						return err
					}

					parsed, err = strconv.ParseFloat(candle.Close, 32)
					if err == nil {
						stockOHLCV.Close = float32(parsed)
					} else {
						return err
					}

					parsedInt, errInt := strconv.ParseInt(candle.Volume, 10, 32)
					if errInt == nil {
						stockOHLCV.Volume = int32(parsedInt)
					} else {
						return errInt
					}

					errTx := tx.Create(&stockOHLCV).Error
					if errTx != nil {
						return errTx
					}
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
