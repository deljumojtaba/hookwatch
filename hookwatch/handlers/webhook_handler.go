package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sajjadgozal/hookwatch/services"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "hookwatch",
		"message": "Service is running",
	})
}

func ReceiveWebhook(c *gin.Context) {
	endpointId := c.Param("endpointId")

	// Extract headers
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		headers[key] = strings.Join(values, ", ")
	}

	// Extract client information
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Parse the webhook data based on method
	var data interface{}

	// For GET requests, try to parse query parameters
	if c.Request.Method == "GET" {
		queryParams := make(map[string]interface{})
		for key, values := range c.Request.URL.Query() {
			if len(values) == 1 {
				queryParams[key] = values[0]
			} else {
				queryParams[key] = values
			}
		}
		data = queryParams
	} else {
		// For other methods, try to parse JSON body
		if err := c.ShouldBindJSON(&data); err != nil {
			// If JSON parsing fails, try to get raw body as string
			if body, err := c.GetRawData(); err == nil && len(body) > 0 {
				data = string(body)
			} else {
				// If no body, create empty object
				data = gin.H{}
			}
		}
	}

	// Process webhook using the service
	webhookService := services.NewWebhookService()
	webhookLog, err := webhookService.ProcessWebhook(
		endpointId,
		c.Request.Method,
		headers,
		data,
		ipAddress,
		userAgent,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process webhook",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Webhook received and processed",
		"endpointId": endpointId,
		"method":     c.Request.Method,
		"webhookId":  webhookLog.ID.Hex(),
		"status":     webhookLog.Status,
		"timestamp":  webhookLog.CreatedAt,
	})
}

func GetWebhookLogs(c *gin.Context) {
	endpointId := c.Param("endpointId")

	// Get limit from query parameter, default to 50
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 50
	}

	// Limit to reasonable range
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 10
	}

	webhookService := services.NewWebhookService()
	webhookLogs, err := webhookService.GetWebhookLogs(endpointId, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve webhook logs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"endpointId": endpointId,
		"logs":       webhookLogs,
		"count":      len(webhookLogs),
		"limit":      limit,
	})
}

func ClearWebhookLogs(c *gin.Context) {
	endpointId := c.Param("endpointId")

	if endpointId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Endpoint ID is required",
		})
		return
	}

	webhookService := services.NewWebhookService()
	deletedCount, err := webhookService.ClearWebhookLogs(endpointId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to clear webhook logs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Webhook logs cleared successfully",
		"endpointId":   endpointId,
		"deletedCount": deletedCount,
	})
}

func ReplayWebhook(c *gin.Context) {
	webhookLogID := c.Param("webhookLogId")

	if webhookLogID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Webhook log ID is required",
		})
		return
	}

	var replayRequest struct {
		TargetURL string `json:"target_url" binding:"required"`
		Timeout   int    `json:"timeout"`
	}

	if err := c.ShouldBindJSON(&replayRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	replayService := services.NewReplayService()
	response, err := replayService.ReplayWebhookByID(webhookLogID, replayRequest.TargetURL, replayRequest.Timeout)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to replay webhook",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Webhook replayed successfully",
		"webhook_log_id": webhookLogID,
		"target_url":     replayRequest.TargetURL,
		"result":         response,
	})
}
