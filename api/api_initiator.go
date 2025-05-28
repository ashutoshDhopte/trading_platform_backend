package api

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func init() {

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
	apiMux.HandleFunc("/buy-stocks", RecoverMiddleware(BuyStocks))
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
				response := getErrorApiResponse("Internal Server Error")
				json.NewEncoder(w).Encode(response)
			}
		}()
		next(w, r)
	}
}
