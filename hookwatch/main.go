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
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Forwarded-Proto", "X-Real-IP"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // Set to false when allowing all origins
	}))

	// Trust proxy headers
	router.SetTrustedProxies([]string{"172.0.0.0/8", "10.0.0.0/8", "192.168.0.0/16"})

	// API routes FIRST (before static files)
	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// Webhook endpoints - accept all HTTP methods
	router.Any("/webhooks/:endpointId/receive", handlers.ReceiveWebhook)
	router.GET("/webhooks/:endpointId/logs", handlers.GetWebhookLogs)
	router.DELETE("/webhooks/:endpointId/logs", handlers.ClearWebhookLogs)
	router.POST("/webhooks/replay/:webhookLogId", handlers.ReplayWebhook)

	// Serve web UI under /ui path to avoid conflicts
	router.Static("/ui", "../web")

	// Redirect root to /ui
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/ui/")
	})

	port := config.GetEnv("PORT", "3000")

	log.Println("Starting HookWatch on port " + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
