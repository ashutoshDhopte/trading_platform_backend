package service

import (
	"trading_platform_backend/db"
	"trading_platform_backend/model"
	"trading_platform_backend/util"
)

func GetMarketData(stockID int64) model.MarketModel {
	// Get specific stock data
	stock := db.GetStockById(stockID)
	if stock.StockID == 0 {
		return model.MarketModel{}
	}

	stockModel := model.StockModel{
		StockID:               stock.StockID,
		Ticker:                stock.Ticker,
		Name:                  stock.Name,
		OpeningPriceDollars:   util.ConvertCentsToDollars(stock.OpeningPriceCents),
		CurrentPriceDollars:   util.ConvertCentsToDollars(stock.CurrentPriceCents),
		UpdatedAt:             util.GetDateTimeString(stock.UpdatedAt),
		OverallSentimentScore: stock.OverallSentimentScore,
	}
	stockModel.ChangedPriceDollars = stockModel.GetChangedPriceDollars()
	stockModel.ChangedPercent = stockModel.GetChangedPercent()

	// Get the 10 most recent news articles for this specific stock
	newsArticles := db.GetLatestNewsArticlesByStock(stockID, 10)
	newsModels := make([]model.NewsModel, 0)

	for _, news := range newsArticles {
		newsModel := model.NewsModel{
			NewsArticleID:   news.NewsArticleID,
			Ticker:          news.Ticker,
			ArticleTitle:    news.ArticleTitle,
			ArticleSummary:  news.ArticleSummary,
			ArticleURL:      news.ArticleURL,
			PublicationTime: util.GetDateTimeString(news.PublicationTime),
			SentimentScore:  news.SentimentScore,
		}
		newsModels = append(newsModels, newsModel)
	}

	return model.MarketModel{
		Stock: stockModel,
		News:  newsModels,
	}
}

func GetAllStocksData() []model.StockModel {
	stocks := db.GetAllStocks()
	stockModels := make([]model.StockModel, 0)

	for _, stock := range stocks {
		stockModel := model.StockModel{
			StockID:             stock.StockID,
			Ticker:              stock.Ticker,
			Name:                stock.Name,
			OpeningPriceDollars: util.ConvertCentsToDollars(stock.OpeningPriceCents),
			CurrentPriceDollars: util.ConvertCentsToDollars(stock.CurrentPriceCents),
			UpdatedAt:           util.GetDateTimeString(stock.UpdatedAt),
		}
		stockModel.ChangedPriceDollars = stockModel.GetChangedPriceDollars()
		stockModel.ChangedPercent = stockModel.GetChangedPercent()

		stockModels = append(stockModels, stockModel)
	}

	return stockModels
}

func GetStockNewsWithPagination(stockID int64, page int) []model.NewsModel {
	// Calculate offset based on page number (page 1 = offset 0, page 2 = offset 10, etc.)
	limit := 10
	offset := (page - 1) * limit

	// Get paginated news articles for this specific stock
	newsArticles := db.GetLatestNewsArticlesByStockWithPagination(stockID, limit, offset)
	newsModels := make([]model.NewsModel, 0)

	for _, news := range newsArticles {
		newsModel := model.NewsModel{
			NewsArticleID:   news.NewsArticleID,
			Ticker:          news.Ticker,
			ArticleTitle:    news.ArticleTitle,
			ArticleSummary:  news.ArticleSummary,
			ArticleURL:      news.ArticleURL,
			PublicationTime: util.GetDateTimeString(news.PublicationTime),
			SentimentScore:  news.SentimentScore,
		}
		newsModels = append(newsModels, newsModel)
	}

	return newsModels
}
