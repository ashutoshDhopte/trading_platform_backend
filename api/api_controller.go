package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"trading_platform_backend/model"
	"trading_platform_backend/service"
)

func GetDashboard(w http.ResponseWriter, r *http.Request) {

	userIdStr := r.URL.Query().Get("userId")
	var response model.ApiResponse
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUserByEmailAndPassword(w http.ResponseWriter, r *http.Request) {

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	var response model.ApiResponse
	if email == "" || password == "" {
		response = getErrorApiResponse("email and password are required")
	} else {
		userModel := service.GetUserByEmailAndPassword(email, password)
		response = getSuccessApiResponse(userModel)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
