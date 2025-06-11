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

func GetUserById(w http.ResponseWriter, r *http.Request) {

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
		if err != nil {
			response = getErrorApiResponse("userId is required")
		} else {
			userModel := service.GetUserById(userId)
			response = getSuccessApiResponse(userModel)
		}
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

func LoginUser(w http.ResponseWriter, r *http.Request) {

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
		return
	}

	result := service.LoginUser(email, password)
	if result > 0 {
		response = getSuccessApiResponse(result)
	} else {
		response = getErrorApiResponse("user not found")
	}
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {

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
	verifyPassword := r.URL.Query().Get("verifyPassword")

	if email == "" || password == "" || verifyPassword == "" {
		response = getErrorApiResponse("email and password are required")
		return
	}

	if password != verifyPassword {
		response = getErrorApiResponse("passwords do not match")
		return
	}

	result, err := service.CreateAccount(email, password)
	if err != nil {
		fmt.Println(err.Error())
		response = getErrorApiResponse("failed to create account, " + err.Error())
	} else if result == 0 {
		response = getErrorApiResponse("failed to create account")
	} else {
		response = getSuccessApiResponse(result)
	}
}

func GetOrders(w http.ResponseWriter, r *http.Request) {

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
		return
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		response = getErrorApiResponse("userId is invalid")
	}

	result := service.GetAllOrders(userId)
	if result != nil {
		response = getSuccessApiResponse(result)
	}
}

func AddStockToWatchlist(w http.ResponseWriter, r *http.Request) {

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
	stockIdStr := r.URL.Query().Get("stockId")
	targetPriceStr := r.URL.Query().Get("targetPrice")

	if userIdStr == "" || stockIdStr == "" || targetPriceStr == "" {
		response = getErrorApiResponse("userId, stockId and targetPrice is required")
		return
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		response = getErrorApiResponse("userId is invalid")
		return
	}

	stockId, err := strconv.ParseInt(stockIdStr, 10, 64)
	if err != nil {
		response = getErrorApiResponse("stockId is invalid")
		return
	}

	targetPrice, err := strconv.ParseFloat(targetPriceStr, 64)
	if err != nil {
		response = getErrorApiResponse("targetPrice is invalid")
		return
	}

	err = service.AddStockToWatchlist(userId, stockId, targetPrice)
	if err != nil {
		response = getErrorApiResponse(err.Error())
	} else {
		response = getSuccessApiResponse("")
	}
}
