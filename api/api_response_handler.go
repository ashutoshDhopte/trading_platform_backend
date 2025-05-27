package api

import (
	"trading_platform_backend/model"
)

func getSuccessApiResponse(data interface{}) model.ApiResponse {
	return model.ApiResponse{Success: true, Data: data}
}

func getErrorApiResponse(errorMessage string) model.ApiResponse {
	return model.ApiResponse{ErrorMessage: errorMessage}
}
