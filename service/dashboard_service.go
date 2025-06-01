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

	user := db.GetUserById(userId)
	if user.UserID == 0 {
		return model.DashboardModel{}
	}

	stocks := db.GetAllStocks()
	var stockModels []model.StockModel
	stockMap := make(map[int64]orm.Stocks)
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
		stockMap[stock.StockID] = stock
	}

	holdings := db.GetActiveHoldingsByUserID(userId)
	var holdingModels []model.HoldingModel
	var totalHoldingValueCents float64
	var totalPnlCents float64
	for _, holding := range holdings {

		var pnl float64
		if holding.Quantity > 0 {
			pnl = stockMap[holding.StockID].OpeningPriceCents - holding.AverageCostPerShareCents
		} else if holding.Quantity < 0 {
			pnl = holding.AverageCostPerShareCents - stockMap[holding.StockID].OpeningPriceCents
		}

		totalPnlCents += pnl

		holdingValueCents := float64(holding.Quantity) * holding.AverageCostPerShareCents
		totalHoldingValueCents += holdingValueCents

		holdingModels = append(holdingModels, model.HoldingModel{
			HoldingID:                  holding.HoldingID,
			StockTicker:                stockMap[holding.StockID].Ticker,
			Quantity:                   holding.Quantity,
			AverageCostPerShareDollars: util.ConvertCentsToDollars(holding.AverageCostPerShareCents),
			TotalValueDollars:          util.ConvertCentsToDollars(holdingValueCents),
			UpdatedAt:                  util.GetDateTimeString(holding.UpdatedAt),
			PnLDollars:                 util.ConvertCentsToDollars(pnl),
			PnLPercent:                 (pnl / holdingValueCents) * 100,
		})
	}

	userModel := model.UserModel{
		UserID:             user.UserID,
		CashBalanceDollars: util.ConvertCentsToDollars(user.CashBalanceCents),
	}

	return model.DashboardModel{
		User:                     userModel,
		Stocks:                   stockModels,
		Holdings:                 holdingModels,
		TotalHoldingValueDollars: util.ConvertCentsToDollars(totalHoldingValueCents),
		PortfolioValueDollars:    util.ConvertCentsToDollars(user.CashBalanceCents + totalHoldingValueCents),
		TotalPnLDollars:          util.ConvertCentsToDollars(totalPnlCents),
		TotalPnLPercent:          (totalPnlCents / totalHoldingValueCents) * 100,
	}
}

func BuyStocks(userId int64, ticker string, quantity int64) string {

	//get stock using ticker
	//if stock is not present, err
	//fetch user
	//if balance is less than total order value, err
	//create new order
	//save the order
	//load the holding using ticker and userid
	//if not exist, create new holding
	//update the values
	//save the holding
	//update the user balance

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		stock := db.GetStockByTicker(ticker)
		if stock.StockID == 0 {
			return errors.New("stock " + ticker + " not found")
		}

		user := db.GetUserById(userId)
		if user.UserID == 0 {
			return errors.New("user does not exist")
		}

		totalOrderValueCents := float64(quantity) * stock.OpeningPriceCents
		if totalOrderValueCents > user.CashBalanceCents {
			return errors.New("user don't have enough balance")
		}

		order := orm.Orders{
			UserID:               userId,
			StockID:              stock.StockID,
			TradeType:            util.TRADE_TYPE_BUY,
			OrderStatus:          util.ORDER_STATUS_EXECUTED,
			Quantity:             quantity,
			PricePerShareCents:   stock.OpeningPriceCents,
			TotalOrderValueCents: float64(quantity) * stock.OpeningPriceCents,
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

		oldTotal := holding.AverageCostPerShareCents * float64(holding.Quantity)
		newQuantity := holding.Quantity + quantity
		holding.AverageCostPerShareCents = (oldTotal + order.TotalOrderValueCents) / float64(newQuantity)

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
