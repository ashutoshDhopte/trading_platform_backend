package db

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB // Global variable to hold the database connection pool

func InitDB() {

	_ = godotenv.Load()

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSL_MODE")

	if dbname == "" || password == "" || host == "" || user == "" || sslmode == "" {
		fmt.Println("Warning: environment variable(s) are not set")
	}
	if port == "" {
		port = "5432" // default fallback
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		user, password, dbname, host, port, sslmode,
	)

	//dbUrl := os.Getenv("DB_URL")
	//
	//if dbUrl == "" {
	//	fmt.Println("Warning: Database URL is not set")
	//}

	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error opening database: %q\n", err)
	} else {
		fmt.Println("Successfully connected to PostgreSQL database!")

		//reset current price
		UpdateStocksResetCurrentPrice()
	}
}
