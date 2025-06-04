package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"trading_platform_backend/model"
	"trading_platform_backend/service"
)

func GetDashboard(w http.ResponseWriter, r *http.Request) {

	var response model.ApiResponse

	//LIFO
	defer func() {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}()

	userIdStr := r.URL.Query().Get("userId")
	if userIdStr == "" {
		response = getErrorApiResponse("userId is required")
	} else {
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err == nil {
			dashboard := service.GetDashboardData(userId)
			response = getSuccessApiResponse(dashboard)
		} else {
			response = getErrorApiResponse("userId is required")
		}
	}
}

func GetUserByEmailAndPassword(w http.ResponseWriter, r *http.Request) {

	var response model.ApiResponse

	//LIFO
	defer func() {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}()

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	if email == "" || password == "" {
		response = getErrorApiResponse("email and password are required")
	} else {
		userModel := service.GetUserByEmailAndPassword(email, password)
		response = getSuccessApiResponse(userModel)
	}
}

func BuyStocks(w http.ResponseWriter, r *http.Request) {

	var response model.ApiResponse

	//LIFO
	defer func() {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}()

	type TradeRequest struct {
		UserID   int64  `json:"userId"`
		Ticker   string `json:"ticker"`
		Quantity int64  `json:"quantity"`
	}

	var payload TradeRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		response = getErrorApiResponse("Invalid payload")
		return
	}

	if payload.UserID == 0 || payload.Quantity == 0 || payload.Ticker == "" {
		response = getErrorApiResponse("Invalid payload")
		return
	}

	result := service.BuyStocks(payload.UserID, payload.Ticker, payload.Quantity)
	if result == "" {
		response = getSuccessApiResponse("")
	} else {
		response = getErrorApiResponse(result)
	}
}

func SellStocks(w http.ResponseWriter, r *http.Request) {

	var response model.ApiResponse

	//LIFO
	defer func() {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}()

	type TradeRequest struct {
		UserID   int64  `json:"userId"`
		Ticker   string `json:"ticker"`
		Quantity int64  `json:"quantity"`
	}

	var payload TradeRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		response = getErrorApiResponse("Invalid payload")
		return
	}

	if payload.UserID == 0 || payload.Quantity == 0 || payload.Ticker == "" {
		response = getErrorApiResponse("Invalid payload")
		return
	}

	result := service.SellStocks(payload.UserID, payload.Ticker, payload.Quantity)
	if result == "" {
		response = getSuccessApiResponse("")
	} else {
		response = getErrorApiResponse(result)
	}
}
