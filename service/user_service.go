package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"trading_platform_backend/auth"
	"trading_platform_backend/db"
	"trading_platform_backend/model"
	"trading_platform_backend/orm"
	"trading_platform_backend/util"
)

func GetUserByEmailAndPassword(email string, password string) model.UserModel {
	user := db.GetUserByEmail(email)

	validPassword := util.CheckPasswordHash(password, user.HashedPassword)
	if !validPassword {
		return model.UserModel{}
	}

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

func LoginUser(email string, password string) (model.Auth, error) {

	user := db.GetUserByEmail(email)

	validPassword := util.CheckPasswordHash(password, user.HashedPassword)
	if !validPassword {
		return model.Auth{}, errors.New("invalid password")
	}

	//jwt session token
	token, err := auth.CreateJWT(int(user.UserID), user.Email)
	if err != nil {
		return model.Auth{}, err
	}

	return model.Auth{
		UserId: int(user.UserID),
		Token:  token,
	}, nil
}

func CreateAccount(email string, password string) (model.Auth, error) {

	var authResponse model.Auth

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		user := db.GetUserByEmail(email)

		if user.UserID > 0 {
			return errors.New("user with this email already exists")
		}

		hashedPassword, err := util.HashPassword(password)
		if err != nil {
			return err
		}

		user = orm.Users{
			Username:         email,
			Email:            email,
			HashedPassword:   hashedPassword,
			CashBalanceCents: util.InitialInvestmentCents,
		}

		err = db.DB.Create(&user).Error
		if err != nil {
			return err
		}

		//jwt session token
		token, err := auth.CreateJWT(int(user.UserID), user.Email)
		if err != nil {
			return err
		}

		authResponse.Token = token
		authResponse.UserId = int(user.UserID)

		return nil
	})

	if err != nil {
		return model.Auth{}, err
	}

	return authResponse, nil
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

func PasswordMigration() {
	var users []orm.Users
	db.DB.Find(&users)
	for i := range users {
		hashedPassword, err := util.HashPassword(users[i].HashedPassword)
		if err == nil {
			users[i].HashedPassword = hashedPassword
		}
	}
	err := db.DB.Save(users).Error
	if err != nil {
		fmt.Println(err)
	}
}
