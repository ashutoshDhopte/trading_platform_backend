package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"trading_platform_backend/db"
	"trading_platform_backend/model"
	"trading_platform_backend/orm"
	"trading_platform_backend/util"
)

const initialInvestmentCents = 10000000

func GetDashboardData(userId int64) model.DashboardModel {

	user := db.GetUserById(userId)
	if user.UserID == 0 {
		return model.DashboardModel{}
	}

	stocks := db.GetAllStocks()
	stockModels := make([]model.StockModel, 0)
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
	holdingModels := make([]model.HoldingModel, 0)
	var totalHoldingValueCents float64
	var totalPnlCents float64
	for _, holding := range holdings {

		var pnl float64
		if holding.Quantity > 0 {
			pnl = stockMap[holding.StockID].OpeningPriceCents - holding.AverageCostPerShareCents
		} else if holding.Quantity < 0 {
			pnl = holding.AverageCostPerShareCents - stockMap[holding.StockID].OpeningPriceCents
		}

		holdingValueCents := float64(holding.Quantity) * holding.AverageCostPerShareCents

		holdingModels = append(holdingModels, model.HoldingModel{
			HoldingID:                  holding.HoldingID,
			StockTicker:                stockMap[holding.StockID].Ticker,
			Quantity:                   holding.Quantity,
			AverageCostPerShareDollars: util.ConvertCentsToDollars(holding.AverageCostPerShareCents),
			TotalValueDollars:          util.ConvertCentsToDollars(holdingValueCents),
			UpdatedAt:                  util.GetDateTimeString(holding.UpdatedAt),
			PnLDollars:                 util.ConvertCentsToDollars(pnl),
			PnLPercent:                 (pnl / math.Abs(holdingValueCents)) * 100,
		})

		totalPnlCents += pnl
		totalHoldingValueCents += holdingValueCents
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
		TotalReturnPercent:       ((user.CashBalanceCents + totalHoldingValueCents - initialInvestmentCents) / initialInvestmentCents) * 100,
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

		holding := db.GetHoldingByUserIdAndStockId(userId, stock.StockID)

		buyQuantity := float64(quantity)
		if holding.HoldingID > 0 && holding.Quantity < 0 {
			buyQuantity = math.Min(math.Abs(float64(holding.Quantity)), float64(quantity)) //to make holding from -ve to 0
		}

		result := buyOrder(tx, &user, stock, int64(buyQuantity), &holding)
		if result != "" {
			return errors.New("Failed to buy stock, " + result)
		}

		//extra quantity for long trade
		if quantity > int64(buyQuantity) {
			longQuantity := quantity - int64(buyQuantity)
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
		PricePerShareCents:   stock.OpeningPriceCents,
		TotalOrderValueCents: float64(quantity) * stock.OpeningPriceCents,
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

	oldTotal := holding.AverageCostPerShareCents * math.Abs(float64(holding.Quantity))
	holding.Quantity += quantity
	if holding.Quantity != 0 {
		holding.AverageCostPerShareCents = math.Abs((oldTotal + order.TotalOrderValueCents) / float64(holding.Quantity))
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

		totalOrderValueCents := float64(quantity) * stock.OpeningPriceCents
		if totalOrderValueCents > user.CashBalanceCents {
			return errors.New("user don't have enough balance")
		}

		holding := db.GetHoldingByUserIdAndStockId(user.UserID, stock.StockID)
		sellQuantity := float64(quantity)
		if holding.HoldingID > 0 && holding.Quantity > 0 {
			sellQuantity = math.Min(math.Abs(float64(holding.Quantity)), float64(quantity)) //to make the holding from +ve to 0
		}

		result := sellOrder(tx, &user, stock, int64(sellQuantity), &holding)
		if result != "" {
			return errors.New("failed to sell order, " + result)
		}

		//extra quantity short trade
		if quantity > int64(sellQuantity) {
			shortQuantity := quantity - int64(sellQuantity)
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
		PricePerShareCents:   stock.OpeningPriceCents,
		TotalOrderValueCents: float64(quantity) * stock.OpeningPriceCents,
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

	oldTotal := holding.AverageCostPerShareCents * math.Abs(float64(holding.Quantity))
	holding.Quantity -= quantity
	if holding.Quantity != 0 {
		holding.AverageCostPerShareCents = math.Abs((oldTotal - order.TotalOrderValueCents) / float64(holding.Quantity))
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
