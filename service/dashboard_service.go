package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"trading_platform_backend/db"
	"trading_platform_backend/model"
	"trading_platform_backend/orm"
	"trading_platform_backend/util"
)

func GetDashboardData(userId int64) model.DashboardModel {

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
		Stocks:   stockModels,
		Holdings: holdingModels,
	}
}

func BuyStocks(userId int64, ticker string, quantity int64) string {

	//get stock using ticker
	//if stock is not present, err
	//create new order
	//save the order
	//load the holding using ticker and userid
	//if not exist, create new holding
	//update the values
	//save the holding

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		stock := db.GetStockByTicker(ticker)
		if stock.StockID == 0 {
			return errors.New("stock " + ticker + " not found")
		}

		user := db.GetUserById(userId)
		if user.UserID == 0 {
			return errors.New("user does not exist")
		}

		totalOrderValueCents := quantity * stock.OpeningPriceCents
		if totalOrderValueCents > user.CashBalanceCents {
			return errors.New("user don't have enough balance")
		}

		order := orm.Orders{
			UserID:               userId,
			StockID:              stock.StockID,
			TradeType:            util.BUY,
			OrderStatus:          util.EXECUTED,
			Quantity:             quantity,
			PricePerShareCents:   stock.OpeningPriceCents,
			TotalOrderValueCents: quantity * stock.OpeningPriceCents,
		}

		if err := db.DB.Create(&order).Error; err != nil {
			fmt.Println("Failed to save order:", err)
			return errors.New("failed to save order")
		}

		holding := db.GetHoldingByUserIdAndStockId(userId, stock.StockID)
		if holding.HoldingID == 0 {
			holding = orm.Holdings{
				StockID: stock.StockID,
				UserID:  userId,
			}
		}

		oldTotal := holding.AverageCostPerShareCents * holding.Quantity
		newQuantity := holding.Quantity + quantity
		holding.AverageCostPerShareCents = (oldTotal + order.TotalOrderValueCents) / newQuantity

		holding.Quantity += quantity

		if err := db.DB.Save(&holding).Error; err != nil {
			fmt.Println("Failed to save holding:", err)
			return errors.New("failed to save holding")
		}

		user.CashBalanceCents -= order.TotalOrderValueCents

		if err := db.DB.Save(&user).Error; err != nil {
			fmt.Println("Failed to save account data:", err)
			return errors.New("failed to save account data")
		}

		return nil
	})

	if err != nil {
		return "Failed to buy stock, " + err.Error()
	}

	return ""
}
