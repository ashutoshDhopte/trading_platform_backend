package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB // Global variable to hold the database connection pool

func init() {
	connStr := "user=trading_user password=trading_password dbname=trading_db sslmode=disable host=localhost port=5432"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error opening database: %q", err)
	}

	err = db.Ping() // Verify the connection
	if err != nil {
		fmt.Println("Error pinging database: %q", err)
	}
	fmt.Println("Successfully connected to PostgreSQL database!")
}
