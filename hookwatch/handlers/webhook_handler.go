package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReceiveWebhook(c *gin.Context) {
	endpointId := c.Param("endpointId")

	// Log the webhook data
	var data interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// TODO: Process webhook data
	// TODO: Store in database
	// TODO: Trigger any configured actions

	c.JSON(http.StatusOK, gin.H{
		"message":    "Webhook received",
		"endpointId": endpointId,
		"data":       data,
	})
}
