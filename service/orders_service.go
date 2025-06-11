package service

import (
	"time"
	"trading_platform_backend/db"
	"trading_platform_backend/model"
	"trading_platform_backend/util"
)

func GetAllOrders(userId int64) []model.OrderModel {

	ordersAndStocks := db.GetOrdersAndStocksByUserId(userId)

	orderModels := make([]model.OrderModel, len(ordersAndStocks))

	for i, order := range ordersAndStocks {
		orderModels[i] = model.OrderModel{
			OrderID:                int64(order["order_id"].(int32)),
			StockTicker:            order["ticker"].(string),
			StockName:              order["name"].(string),
			TradeType:              order["trade_type"].(string),
			OrderStatus:            order["order_status"].(string),
			Quantity:               order["quantity"].(int64),
			PricePerShareDollars:   util.ConvertCentsToDollars(order["price_per_share_cents"].(int64)),
			TotalOrderValueDollars: util.ConvertCentsToDollars(order["total_order_value_cents"].(int64)),
			CreatedAt:              util.GetDateTimeString(order["created_at"].(time.Time)),
			Notes:                  order["notes"].(string),
		}
	}

	return orderModels
}
