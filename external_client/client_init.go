package external_client

import "github.com/Finnhub-Stock-API/finnhub-go/v2"

type Client struct {
	finnhubClient *finnhub.DefaultApiService
}

var client = Client{}

func InitExternalClient() {
	initFinnhubClient()
}
