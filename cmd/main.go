package main

import (
	"auth-session/config"
	"auth-session/router"
	"log"
)

func main() {
	db := config.ConnectDB()
	r := router.SetupRouter(db)

	log.Println("🚀 Server running at http://localhost:8080")
	r.Run(":8080")
}
