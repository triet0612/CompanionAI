package main

import (
	"CompanionBackend/pkg/config"
	"CompanionBackend/pkg/db"
	"CompanionBackend/pkg/server"
)

func main() {
	start()
}

func start() {
	config := config.Init()
	db := db.Init(config)

	server := server.Init(db, config)
	server.Run()
}
