package api

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
	"runtime/debug"
	"trading_platform_backend/routine"
)

func InitApi() {

	apiMux := registerRoutes()

	mux := http.NewServeMux()
	mux.Handle("/trade-sim/", http.StripPrefix("/trade-sim", apiMux))

	// Setup CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)

	port := ":8080"
	fmt.Println("Server running at http://localhost" + port + "/trade-sim/")
	if err := http.ListenAndServe(port, corsHandler); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

// Factory function to register all routes
func registerRoutes() *http.ServeMux {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/", homeHandler)
	apiMux.HandleFunc("/dashboard", RecoverMiddleware(GetDashboard))
	apiMux.HandleFunc("/user", RecoverMiddleware(GetUserByEmailAndPassword))
	apiMux.HandleFunc("/user/v2", RecoverMiddleware(GetUserById))
	apiMux.HandleFunc("/buy-stocks", RecoverMiddleware(BuyStocks))
	apiMux.HandleFunc("/sell-stocks", RecoverMiddleware(SellStocks))
	apiMux.HandleFunc("/login", RecoverMiddleware(LoginUser))
	apiMux.HandleFunc("/create-account", RecoverMiddleware(CreateAccount))
	apiMux.HandleFunc("/orders", RecoverMiddleware(GetOrders))
	apiMux.HandleFunc("/add-stock-watchlist", RecoverMiddleware(AddStockToWatchlist))
	apiMux.HandleFunc("/delete-stock-watchlist", RecoverMiddleware(DeleteStockFromWatchlist))
	apiMux.HandleFunc("/update-user-setting", RecoverMiddleware(UpdateUserSettings))

	apiMux.HandleFunc("/ws/dashboard", routine.ServeWs)
	// Add more handlers here

	return apiMux
}

// Handler for the root route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// A simple JSON response
	response := map[string]string{"status": "Go backend is operational"}
	json.NewEncoder(w).Encode(response)
}

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
