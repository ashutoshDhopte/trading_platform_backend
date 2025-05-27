package service

import (
	"trading_platform_backend/db"
	"trading_platform_backend/model"
	"trading_platform_backend/orm"
)

func GetDashboardData(userId int64) model.Dashboard {

	user := db.GetUserById(userId)

	userInfo := orm.Users{
		UserID:           userId,
		Username:         user.Username,
		Email:            user.Email,
		CashBalanceCents: user.CashBalanceCents,
		CreatedAt:        user.CreatedAt,
	}

	stocks := db.GetAllStocks()
	holdings := db.GetActiveHoldingsByUserID(userId)

	return model.Dashboard{
		User:     userInfo,
		Stocks:   stocks,
		Holdings: holdings,
	}
}
