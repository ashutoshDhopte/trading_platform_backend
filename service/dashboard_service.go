package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"time"
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
	stockModels := make([]model.StockModel, 0)
	stockMap := make(map[int64]orm.Stocks)
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
		stockMap[stock.StockID] = stock
	}

	holdings := db.GetActiveHoldingsByUserID(userId)
	holdingModels := make([]model.HoldingModel, 0)

	var totalHoldingValueCents int64
	var totalPnlCents int64

	for _, holding := range holdings {

		var pnlCents int64
		if holding.Quantity > 0 {
			pnlCents = stockMap[holding.StockID].CurrentPriceCents - holding.AverageCostPerShareCents
		} else if holding.Quantity < 0 {
			pnlCents = holding.AverageCostPerShareCents - stockMap[holding.StockID].CurrentPriceCents
		}

		holdingValueCents := holding.Quantity * holding.AverageCostPerShareCents

		holdingModels = append(holdingModels, model.HoldingModel{
			HoldingID:                  holding.HoldingID,
			StockTicker:                stockMap[holding.StockID].Ticker,
			Quantity:                   holding.Quantity,
			AverageCostPerShareDollars: util.ConvertCentsToDollars(holding.AverageCostPerShareCents),
			TotalValueDollars:          util.ConvertCentsToDollars(holdingValueCents),
			UpdatedAt:                  util.GetDateTimeString(holding.UpdatedAt),
			PnLDollars:                 util.ConvertCentsToDollars(pnlCents),
			PnLPercent:                 (float64(pnlCents) / math.Abs(float64(holdingValueCents))) * 100,
		})

		totalPnlCents += pnlCents
		totalHoldingValueCents += holdingValueCents
	}

	userModel := model.UserModel{
		UserID:             user.UserID,
		Username:           user.Username,
		Email:              user.Email,
		CashBalanceDollars: util.ConvertCentsToDollars(user.CashBalanceCents),
		CreatedAt:          util.GetDateTimeString(user.CreatedAt),
		UpdatedAt:          util.GetDateTimeString(user.UpdatedAt),
	}

	return model.DashboardModel{
		User:                     userModel,
		Stocks:                   stockModels,
		Holdings:                 holdingModels,
		TotalHoldingValueDollars: util.ConvertCentsToDollars(totalHoldingValueCents),
		PortfolioValueDollars:    util.ConvertCentsToDollars(user.CashBalanceCents + totalHoldingValueCents),
		TotalPnLDollars:          util.ConvertCentsToDollars(totalPnlCents),
		TotalReturnPercent:       (float64(user.CashBalanceCents+totalHoldingValueCents-util.InitialInvestmentCents) / util.InitialInvestmentCents) * 100,
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

		totalOrderValueCents := quantity * stock.CurrentPriceCents
		if totalOrderValueCents > user.CashBalanceCents {
			return errors.New("user don't have enough balance")
		}

		holding := db.GetHoldingByUserIdAndStockId(userId, stock.StockID)

		buyQuantity := quantity
		if holding.HoldingID > 0 && holding.Quantity < 0 {
			buyQuantity = int64(math.Min(math.Abs(float64(holding.Quantity)), float64(quantity))) //to make holding from -ve to 0
		}

		result := buyOrder(tx, &user, stock, buyQuantity, &holding)
		if result != "" {
			return errors.New("Failed to buy stock, " + result)
		}

		//extra quantity for long trade
		if quantity > buyQuantity {
			longQuantity := quantity - buyQuantity
			result = buyOrder(tx, &user, stock, longQuantity, &holding)
			if result != "" {
				return errors.New("Failed to buy stock, " + result)
			}
		}

		return nil
	})

	if err != nil {
		return "Failed to buy stock, " + err.Error()
	}

	return ""
}

func buyOrder(tx *gorm.DB, user *orm.Users, stock orm.Stocks, quantity int64, holding *orm.Holdings) string {

	order := orm.Orders{
		UserID:               user.UserID,
		StockID:              stock.StockID,
		TradeType:            util.TradeTypeBuy,
		OrderStatus:          util.OrderStatusExecuted,
		Quantity:             quantity,
		PricePerShareCents:   stock.CurrentPriceCents,
		TotalOrderValueCents: quantity * stock.CurrentPriceCents,
		CreatedAt:            time.Now(),
	}

	if err := tx.Create(&order).Error; err != nil {
		fmt.Println("Failed to save order:", err)
		return "failed to save order"
	}

	if holding.HoldingID == 0 {
		*holding = orm.Holdings{
			StockID: stock.StockID,
			UserID:  user.UserID,
		}
	}

	oldTotalCents := holding.AverageCostPerShareCents * int64(math.Abs(float64(holding.Quantity)))
	holding.Quantity += quantity
	if holding.Quantity != 0 {
		holding.AverageCostPerShareCents = int64(math.Abs(float64((oldTotalCents + order.TotalOrderValueCents) / holding.Quantity)))
	}

	if err := tx.Save(&holding).Error; err != nil {
		fmt.Println("Failed to save holding:", err)
		return "failed to save holding"
	}

	user.CashBalanceCents -= order.TotalOrderValueCents

	if err := tx.Save(&user).Error; err != nil {
		fmt.Println("Failed to save account data:", err)
		return "failed to save account data"
	}

	return ""
}

func SellStocks(userId int64, ticker string, quantity int64) string {

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

		totalOrderValueCents := quantity * stock.CurrentPriceCents
		if totalOrderValueCents > user.CashBalanceCents {
			return errors.New("user don't have enough balance")
		}

		holding := db.GetHoldingByUserIdAndStockId(user.UserID, stock.StockID)
		sellQuantity := quantity
		if holding.HoldingID > 0 && holding.Quantity > 0 {
			sellQuantity = int64(math.Min(math.Abs(float64(holding.Quantity)), float64(quantity))) //to make the holding from +ve to 0
		}

		result := sellOrder(tx, &user, stock, sellQuantity, &holding)
		if result != "" {
			return errors.New("failed to sell order, " + result)
		}

		//extra quantity short trade
		if quantity > sellQuantity {
			shortQuantity := quantity - sellQuantity
			result = sellOrder(tx, &user, stock, shortQuantity, &holding)
			if result != "" {
				return errors.New("failed to sell order, " + result)
			}
		}

		return nil
	})

	if err != nil {
		return "Failed to sell stock, " + err.Error()
	}

	return ""
}

func sellOrder(tx *gorm.DB, user *orm.Users, stock orm.Stocks, quantity int64, holding *orm.Holdings) string {

	order := orm.Orders{
		UserID:               user.UserID,
		StockID:              stock.StockID,
		TradeType:            util.TradeTypeSell,
		OrderStatus:          util.OrderStatusExecuted,
		Quantity:             quantity,
		PricePerShareCents:   stock.CurrentPriceCents,
		TotalOrderValueCents: quantity * stock.CurrentPriceCents,
		CreatedAt:            time.Now(),
	}

	if err := tx.Create(&order).Error; err != nil {
		fmt.Println("Failed to save order:", err)
		return "failed to save order"
	}

	if holding.HoldingID == 0 {
		*holding = orm.Holdings{
			StockID: stock.StockID,
			UserID:  user.UserID,
		}
	}

	oldTotalCents := holding.AverageCostPerShareCents * int64(math.Abs(float64(holding.Quantity)))
	holding.Quantity -= quantity
	if holding.Quantity != 0 {
		holding.AverageCostPerShareCents = int64(math.Abs(float64((oldTotalCents - order.TotalOrderValueCents) / holding.Quantity)))
	}

	if err := tx.Save(&holding).Error; err != nil {
		fmt.Println("Failed to save holding:", err)
		return "failed to save holding"
	}

	user.CashBalanceCents += order.TotalOrderValueCents

	if err := tx.Save(&user).Error; err != nil {
		fmt.Println("Failed to save account data:", err)
		return "failed to save account data"
	}

	return ""
}
