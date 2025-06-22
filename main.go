package main

import (
	"trading_platform_backend/api"
	"trading_platform_backend/db"
	"trading_platform_backend/external_client"
	"trading_platform_backend/routine"
)

func main() {
	initPackages()
}

func initPackages() {
	//maintain sequence
	db.InitDB()
	routine.InitRoutines()
	external_client.InitExternalClient()
	//========= add new init before this
	api.InitApi()
}
