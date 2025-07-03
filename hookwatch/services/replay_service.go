package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sajjadgozal/hookwatch/models"
)

type ReplayService struct{}

func NewReplayService() *ReplayService {
	return &ReplayService{}
}

type ReplayRequest struct {
	TargetURL    string            `json:"target_url"`
	Method       string            `json:"method"`
	Headers      map[string]string `json:"headers"`
	Body         interface{}       `json:"body"`
	Timeout      int               `json:"timeout"` // in seconds
	WebhookLogID string            `json:"webhook_log_id"`
}

type ReplayResponse struct {
	Success      bool              `json:"success"`
	StatusCode   int               `json:"status_code"`
	ResponseBody string            `json:"response_body"`
	Headers      map[string]string `json:"response_headers"`
	Duration     time.Duration     `json:"duration"`
	Error        string            `json:"error,omitempty"`
}

// ReplayWebhook sends a webhook log to an external endpoint
func (s *ReplayService) ReplayWebhook(webhookLog *models.WebhookLog, targetURL string, timeout int) (*ReplayResponse, error) {
	startTime := time.Now()

	// Set default timeout if not provided
	if timeout <= 0 {
		timeout = 30
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Prepare request body
	var body io.Reader
	if webhookLog.Body != nil {
		jsonBody, err := json.Marshal(webhookLog.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %v", err)
		}
		body = bytes.NewBuffer(jsonBody)
	}

	// Create HTTP request
	req, err := http.NewRequest(webhookLog.Method, targetURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers from original webhook
	for key, value := range webhookLog.Headers {
		req.Header.Set(key, value)
	}

	// Ensure Content-Type is set for JSON requests
	if webhookLog.Body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return &ReplayResponse{
			Success:  false,
			Error:    err.Error(),
			Duration: time.Since(startTime),
		}, nil
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ReplayResponse{
			Success:  false,
			Error:    fmt.Sprintf("failed to read response: %v", err),
			Duration: time.Since(startTime),
		}, nil
	}

	// Extract response headers
	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			respHeaders[key] = values[0]
		}
	}

	// Log the replay attempt
	log.Printf("Replayed webhook %s to %s - Status: %d, Duration: %v",
		webhookLog.ID.Hex(), targetURL, resp.StatusCode, time.Since(startTime))

	return &ReplayResponse{
		Success:      resp.StatusCode >= 200 && resp.StatusCode < 300,
		StatusCode:   resp.StatusCode,
		ResponseBody: string(respBody),
		Headers:      respHeaders,
		Duration:     time.Since(startTime),
	}, nil
}

// ReplayWebhookByID replays a webhook log by its ID
func (s *ReplayService) ReplayWebhookByID(webhookLogID string, targetURL string, timeout int) (*ReplayResponse, error) {
	// Get webhook log from database
	webhookService := NewWebhookService()
	webhookLog, err := webhookService.GetWebhookLogByID(webhookLogID)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook log: %v", err)
	}

	return s.ReplayWebhook(webhookLog, targetURL, timeout)
}
