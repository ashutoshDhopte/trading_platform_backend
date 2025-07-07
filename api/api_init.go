package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"trading_platform_backend/routine"

	"github.com/rs/cors"
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

	apiMux.HandleFunc("/login", RecoverMiddleware(LoginUser))
	apiMux.HandleFunc("/create-account", RecoverMiddleware(CreateAccount))

	apiMux.HandleFunc("/dashboard", JwtMiddleware(GetDashboard))
	apiMux.HandleFunc("/stocks", JwtMiddleware(GetAllStocks))
	apiMux.HandleFunc("/stock-news", GetStockNews)
	apiMux.HandleFunc("/user", JwtMiddleware(GetUserByEmailAndPassword))
	apiMux.HandleFunc("/user/v2", JwtMiddleware(GetUserById))
	apiMux.HandleFunc("/buy-stocks", JwtMiddleware(BuyStocks))
	apiMux.HandleFunc("/sell-stocks", JwtMiddleware(SellStocks))
	apiMux.HandleFunc("/orders", JwtMiddleware(GetOrders))
	apiMux.HandleFunc("/add-stock-watchlist", JwtMiddleware(AddStockToWatchlist))
	apiMux.HandleFunc("/delete-stock-watchlist", JwtMiddleware(DeleteStockFromWatchlist))
	apiMux.HandleFunc("/update-user-setting", JwtMiddleware(UpdateUserSettings))

	apiMux.HandleFunc("/ws/dashboard", RecoverMiddleware(routine.ServeWs))
	apiMux.HandleFunc("/ws/market", RecoverMiddleware(routine.ServeMarketWs))

	//apiMux.HandleFunc("/migrate-passwords", RecoverMiddleware(PasswordMigration))
	//apiMux.HandleFunc("/migrate-news", RecoverMiddleware(NewsMigration))
	apiMux.HandleFunc("/migrate-stock-ohlcv", RecoverMiddleware(MigrateStockOHLCV))
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
