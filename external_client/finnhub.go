package external_client

import (
	"context"
	"log"
	"os"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/joho/godotenv"
)

func initFinnhubClient() {
	_ = godotenv.Load()
	finnhubToken := os.Getenv("FINNHUB_TOKEN")
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", finnhubToken)
	client.finnhubClient = finnhub.NewAPIClient(cfg).DefaultApi
}

func FetchNewsFromFinnhub(ticker string, from string, to string) []finnhub.CompanyNews {

	res, _, err := client.finnhubClient.
		CompanyNews(context.Background()).
		Symbol(ticker).
		From(from).
		To(to).
		Execute()

	if err != nil {
		log.Fatal(err)
		return []finnhub.CompanyNews{}
	}

	companyNewsList := make([]finnhub.CompanyNews, 0)
	for _, companyNews := range res {
		if companyNews.GetHeadline() == "" || companyNews.GetSummary() == "" {
			continue
		}
		companyNewsList = append(companyNewsList, companyNews)
	}

	return companyNewsList
}
