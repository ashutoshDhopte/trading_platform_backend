package db

import (
	"trading_platform_backend/orm"
)

func GetAllStocks() []orm.Stocks {
	var stocks []orm.Stocks
	DB.Find(&stocks)
	return stocks
}

func GetUserByEmailAndPassword(email string, password string) orm.Users {
	var user orm.Users
	DB.Where("email = ? and hashed_password = ?", email, password).First(&user)
	return user
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
