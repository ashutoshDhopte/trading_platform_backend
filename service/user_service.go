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
		NotificationsOn:    user.NotificationsOn,
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
		NotificationsOn:    user.NotificationsOn,
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

func UpdateUserSettings(userId int64, userSettings map[string]interface{}) (model.UserModel, error) {

	userModel := model.UserModel{}

	user := db.GetUserById(userId)
	if user.UserID == 0 {
		return userModel, errors.New("user not found")
	}

	for key, value := range userSettings {
		switch key {
		case "notifications":
			user.NotificationsOn = value.(bool)
		}
	}

	userModel.UserID = userId
	userModel.NotificationsOn = user.NotificationsOn
	userModel.Email = user.Email
	userModel.Username = user.Username
	userModel.CashBalanceDollars = util.ConvertCentsToDollars(user.CashBalanceCents)
	userModel.CreatedAt = util.GetDateTimeString(user.CreatedAt)
	userModel.UpdatedAt = util.GetDateTimeString(user.UpdatedAt)

	return userModel, db.DB.Save(user).Error
}
