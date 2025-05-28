package service

import (
	"trading_platform_backend/db"
	"trading_platform_backend/model"
	"trading_platform_backend/util"
)

func GetUserByEmailAndPassword(email string, password string) model.UserModel {
	user := db.GetUserByEmailAndPassword(email, password)
	return model.UserModel{
		UserID:             user.UserID,
		Username:           user.Username,
		Email:              user.Email,
		CashBalanceDollars: util.ConvertCentsToDollars(user.CashBalanceCents),
		CreatedAt:          util.GetDateTimeString(user.CreatedAt),
		UpdatedAt:          util.GetDateTimeString(user.UpdatedAt),
	}
}
