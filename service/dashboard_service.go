package service

import (
	"trading_platform_backend/db"
	"trading_platform_backend/model"
	"trading_platform_backend/util"
)

func GetDashboardData(userId int64) model.DashboardModel {

	user := db.GetUserById(userId)

	userModel := model.UserModel{
		UserID:             userId,
		Username:           user.Username,
		Email:              user.Email,
		CashBalanceDollars: util.ConvertCentsToDollars(user.CashBalanceCents),
		CreatedAt:          util.GetDateTimeString(user.CreatedAt),
	}

	stocks := db.GetAllStocks()
	var stockModels []model.StockModel
	stockIdTickerMap := make(map[int64]string)
	for _, stock := range stocks {
		stockModels = append(stockModels, model.StockModel{
			StockID:             stock.StockID,
			Ticker:              stock.Ticker,
			Name:                stock.Name,
			OpeningPriceDollars: util.ConvertCentsToDollars(stock.OpeningPriceCents),
			ChangedPriceDollars: 0,
			ChangedPercent:      0,
			UpdatedAt:           util.GetDateTimeString(stock.UpdatedAt),
		})
		stockIdTickerMap[stock.StockID] = stock.Ticker
	}

	holdings := db.GetActiveHoldingsByUserID(userId)
	var holdingModels []model.HoldingModel
	for _, holding := range holdings {
		holdingModels = append(holdingModels, model.HoldingModel{
			HoldingID:                  holding.HoldingID,
			StockTicker:                stockIdTickerMap[holding.StockID],
			Quantity:                   holding.Quantity,
			AverageCostPerShareDollars: util.ConvertCentsToDollars(holding.AverageCostPerShareCents),
			TotalValueDollars:          util.ConvertCentsToDollars(holding.Quantity * holding.AverageCostPerShareCents),
			UpdatedAt:                  util.GetDateTimeString(holding.UpdatedAt),
		})
	}

	return model.DashboardModel{
		User:     userModel,
		Stocks:   stockModels,
		Holdings: holdingModels,
	}
}
