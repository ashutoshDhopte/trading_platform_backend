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

	type AddWatchlistRequest struct {
		UserId      int32   `json:"userId"`
		StockId     int32   `json:"stockId"`
		TargetPrice float64 `json:"targetPrice"`
	}

	var payload AddWatchlistRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		response = getErrorApiResponse("Invalid payload")
		return
	}

	if payload.UserId == 0 || payload.StockId == 0 || payload.TargetPrice <= 0 {
		response = getErrorApiResponse("Invalid payload")
		return
	}

	err = service.AddStockToWatchlist(payload.UserId, payload.StockId, payload.TargetPrice)
	if err != nil {
		response = getErrorApiResponse(err.Error())
	} else {
		response = getSuccessApiResponse("")
	}
}

func DeleteStockFromWatchlist(w http.ResponseWriter, r *http.Request) {

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

	if userIdStr == "" || stockIdStr == "" {
		response = getErrorApiResponse("userId and stockId are required")
		return
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		response = getErrorApiResponse("userId is invalid")
	}

	stockId, err := strconv.ParseInt(stockIdStr, 10, 32)
	if err != nil {
		response = getErrorApiResponse("stockId is invalid")
	}

	err = service.DeleteFromWatchlist(int32(userId), int32(stockId))
	if err != nil {
		response = getErrorApiResponse(err.Error())
	} else {
		response = getSuccessApiResponse("")
	}
}

func UpdateUserSettings(w http.ResponseWriter, r *http.Request) {

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

	var payload map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		response = getErrorApiResponse("Invalid payload")
		return
	}

	userId := payload["userId"].(float64)
	if err != nil {
		response = getErrorApiResponse("userId is invalid")
		return
	}

	settings := payload["settings"].(map[string]interface{})
	if len(settings) == 0 {
		fmt.Println("settings is required")
		response = getErrorApiResponse("settings are empty")
		return
	}

	user, resRrr := service.UpdateUserSettings(int64(userId), settings)
	if resRrr != nil {
		fmt.Println(resRrr)
		response = getErrorApiResponse(resRrr.Error())
	} else {
		response = getSuccessApiResponse(user)
	}
}
