package service

import (
	"errors"
	"gorm.io/gorm"
	"trading_platform_backend/db"
	"trading_platform_backend/model"
	"trading_platform_backend/orm"
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

func GetUserById(userId int64) model.UserModel {
	user := db.GetUserById(userId)
	return model.UserModel{
		UserID:             user.UserID,
		Username:           user.Username,
		Email:              user.Email,
		CashBalanceDollars: util.ConvertCentsToDollars(user.CashBalanceCents),
		CreatedAt:          util.GetDateTimeString(user.CreatedAt),
		UpdatedAt:          util.GetDateTimeString(user.UpdatedAt),
	}
}

func LoginUser(email string, password string) int64 {

	user := db.GetUserByEmailAndPassword(email, password)

	return user.UserID
}

func CreateAccount(email string, password string) (int64, error) {

	var userId int64

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		user := db.GetUserByEmail(email)

		if user.UserID > 0 {
			return errors.New("user with this email already exists")
		}

		user = orm.Users{
			Username:         email,
			Email:            email,
			HashedPassword:   password,
			CashBalanceCents: util.InitialInvestmentCents,
		}

		err := db.DB.Create(&user).Error
		if err != nil {
			return err
		}

		userId = user.UserID

		return nil
	})

	if err != nil {
		return 0, err
	}

	return userId, nil
}
