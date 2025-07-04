package db

import (
	"fmt"
	"time"
	"trading_platform_backend/orm"

	"gorm.io/gorm"
)

func GetAllStocks() []orm.Stocks {
	var stocks []orm.Stocks
	DB.Find(&stocks)
	return stocks
}

func GetUserByEmail(email string) orm.Users {
	var user orm.Users
	DB.Where("email = ?", email).First(&user)
	return user
}

func GetUserById(userId int64) orm.Users {
	var user orm.Users
	DB.Find(&user, userId)
	return user
}

func GetActiveHoldingsByUserID(userID int64) []orm.Holdings {
	var holdings []orm.Holdings
	DB.Where("user_id = ? and quantity != 0", userID).Find(&holdings)
	return holdings
}

func GetStockByTicker(ticker string) orm.Stocks {
	var stock orm.Stocks
	DB.Where("ticker = ?", ticker).First(&stock)
	return stock
}

func GetHoldingByUserIdAndStockId(userId int64, stockId int64) orm.Holdings {
	var holding orm.Holdings
	DB.Where("user_id = ? and stock_id = ?", userId, stockId).First(&holding)
	return holding
}

func UpdateStocksResetCurrentPrice() {
	result := DB.Model(&orm.Stocks{}).Where("1 = 1").Update("current_price_cents", gorm.Expr("opening_price_cents"))
	if result.Error != nil {
		fmt.Println("Failed to reset current stock price, " + result.Error.Error())
	}
}

func GetOrdersByUserId(userId int64) []orm.Orders {
	var orders []orm.Orders
	DB.Where("user_id = ?", userId).Find(&orders)
	return orders
}

func GetOrdersAndStocksByUserId(userId int64) []map[string]interface{} {
	var result []map[string]interface{}
	DB.Table("orders").
		Select("*").
		Joins("join stocks on stocks.stock_id = orders.stock_id").
		Where("orders.user_id = ?", userId).
		Order("orders.created_at desc").
		Find(&result)
	return result
}

func GetStockWatchlistByUserIdAndStockId(userId int32, stockId int32) orm.StockWatchlist {
	var stockWatchlist orm.StockWatchlist
	DB.Where("user_id = ? and stock_id = ? and is_active = true", userId, stockId).Find(&stockWatchlist)
	return stockWatchlist
}

func GetStockWatchlistByUserId(userId int32) []orm.StockWatchlist {
	var stockWatchlist []orm.StockWatchlist
	DB.Where("user_id = ? and is_active = true", userId).Find(&stockWatchlist)
	return stockWatchlist
}

func GetStockWatchlistAndStockTickerByUserId(userId int64) []map[string]interface{} {
	var stockWatchlist []map[string]interface{}
	DB.Table("stock_watchlist").
		Select("stock_watchlist.*, stocks.ticker, stocks.name").
		Joins("join stocks on stocks.stock_id = stock_watchlist.stock_id").
		Where("stock_watchlist.user_id = ? and is_active = true", userId).
		Find(&stockWatchlist)
	return stockWatchlist
}

func GetLatestNewsArticles(limit int) []orm.NewsArticles {
	var newsArticles []orm.NewsArticles
	DB.Order("publication_time desc").Limit(limit).Find(&newsArticles)
	return newsArticles
}

func GetStockById(stockId int64) orm.Stocks {
	var stock orm.Stocks
	DB.Find(&stock, stockId)
	return stock
}

func GetLatestNewsArticlesByStock(stockId int64, limit int) []orm.NewsArticles {
	var newsArticles []orm.NewsArticles
	stock := GetStockById(stockId)
	if stock.StockID == 0 {
		return newsArticles
	}
	DB.Where("ticker = ?", stock.Ticker).Order("publication_time desc").Limit(limit).Find(&newsArticles)
	return newsArticles
}

func GetLatestNewsArticlesByStockWithPagination(stockId int64, limit int, offset int) []orm.NewsArticles {
	var newsArticles []orm.NewsArticles
	stock := GetStockById(stockId)
	if stock.StockID == 0 {
		return newsArticles
	}
	DB.Where("ticker = ?", stock.Ticker).Order("publication_time desc").Limit(limit).Offset(offset).Find(&newsArticles)
	return newsArticles
}

func GetSentimentDataForStock(stockId int64, days int) []orm.NewsArticles {
	var newsArticles []orm.NewsArticles
	stock := GetStockById(stockId)
	if stock.StockID == 0 {
		return newsArticles
	}

	// Calculate the date from 'days' ago
	cutoffDate := time.Now().AddDate(0, 0, -days)

	DB.Where("ticker = ? AND publication_time >= ?", stock.Ticker, cutoffDate).
		Order("publication_time asc").
		Find(&newsArticles)
	return newsArticles
}
