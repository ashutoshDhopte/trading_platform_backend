package db

import "trading_platform_backend/orm"

func GetAllStocks() []orm.Stocks {
	var stocks []orm.Stocks
	DB.Find(&stocks)
	return stocks
}

func GetUserByEmailAndPassword(email string, password string) orm.Users {
	var user orm.Users
	DB.Where("email = ? and hashedPassword = ?", email, password).First(&user)
	return user
}

func GetUserById(userId int64) orm.Users {
	var user orm.Users
	DB.Find(&user, userId)
	return user
}

func GetActiveHoldingsByUserID(userID int64) []orm.Holdings {
	var holdings []orm.Holdings
	DB.Where("user_id = ? and quantity > 0", userID).Find(&holdings)
	return holdings
}
