package main

import (
	"log"

	"github.com/sajjadgozal/hookwatch/config"
	"github.com/sajjadgozal/hookwatch/db"
	"github.com/sajjadgozal/hookwatch/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db.InitMongo(config.GetEnv("MONGO_URI", "mongodb://localhost:27017"))

	router := gin.Default()

	router.POST("/webhooks/:endpointId/receive", handlers.ReceiveWebhook)

	port := config.GetEnv("PORT", "3000")

	log.Println("Starting HookWatch on port " + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
