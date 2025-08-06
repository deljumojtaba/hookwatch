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
		AllowOrigins:     []string{"*"}, // Allow all origins for production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // Set to false when allowing all origins
	}))

	// API routes FIRST (before static files)
	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// Webhook endpoints - accept all HTTP methods
	router.Any("/webhooks/:endpointId/receive", handlers.ReceiveWebhook)
	router.GET("/webhooks/:endpointId/logs", handlers.GetWebhookLogs)
	router.DELETE("/webhooks/:endpointId/logs", handlers.ClearWebhookLogs)
	router.POST("/webhooks/replay/:webhookLogId", handlers.ReplayWebhook)

	// Static files LAST (wildcard route)
	router.Static("/", "../web")

	port := config.GetEnv("PORT", "3000")

	log.Println("Starting HookWatch on port " + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
