package db

import (
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // Global variable to hold the database connection pool

func init() {
	connStr := "user=trading_user password=trading_password dbname=trading_db sslmode=disable host=localhost port=5432"
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		fmt.Println("Error opening database: %q", err)
	}
	fmt.Println("Successfully connected to PostgreSQL database!")
}
