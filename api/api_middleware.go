package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"trading_platform_backend/auth"
)

func RecoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				fmt.Println("Recovered from panic:", rec)
				fmt.Printf("%s\n", debug.Stack())
				response := getErrorApiResponse("Internal Server Error")
				json.NewEncoder(w).Encode(response)
			}
		}()
		next(w, r)
	}
}

func JwtMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if rec := recover(); rec != nil {
				fmt.Println("Recovered from panic:", rec)
				fmt.Printf("%s\n", debug.Stack())
				response := getErrorApiResponse("Internal Server Error")
				json.NewEncoder(w).Encode(response)
			}
		}()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		_, err := auth.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
