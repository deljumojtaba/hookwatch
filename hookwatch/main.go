package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sajjadgozal/hookwatch/config"
	"github.com/sajjadgozal/hookwatch/db"
	"github.com/sajjadgozal/hookwatch/handlers"
)

func main() {
	config.LoadEnv()

	db.InitMongo(config.GetEnv("MONGO_URI", "mongodb://localhost:27017"))

	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	router.GET("/", handlers.HealthCheck)

	// Webhook endpoints - accept all HTTP methods
	router.Any("/webhooks/:endpointId/receive", handlers.ReceiveWebhook)
	router.GET("/webhooks/:endpointId/logs", handlers.GetWebhookLogs)
	router.DELETE("/webhooks/:endpointId/logs", handlers.ClearWebhookLogs)
	router.POST("/webhooks/replay/:webhookLogId", handlers.ReplayWebhook)

	port := config.GetEnv("PORT", "3000")

	log.Println("Starting HookWatch on port " + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
