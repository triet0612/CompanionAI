package main

import (
	"CompanionBackend/pkg/api"
	"CompanionBackend/pkg/config"
	"CompanionBackend/pkg/db"
)

func main() {
	start()
}

func start() {
	config := config.Init()
	db := db.Init(config)

	server := api.Init(db, config)
	server.Run()
}
