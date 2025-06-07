package main

import (
	"trading_platform_backend/api"
	"trading_platform_backend/db"
	"trading_platform_backend/routine"
)

func main() {
	initPackages()
}

func initPackages() {
	//maintain sequence
	db.InitDB()
	routine.InitRoutines()
	api.InitApi()
}
