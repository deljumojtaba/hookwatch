package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sajjadgozal/hookwatch/db"
	"github.com/sajjadgozal/hookwatch/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WebhookService struct{}

func NewWebhookService() *WebhookService {
	return &WebhookService{}
}

// ProcessWebhook handles the complete webhook processing pipeline
func (s *WebhookService) ProcessWebhook(endpointID string, method string, headers map[string]string, body interface{}, ipAddress, userAgent string) (*models.WebhookLog, error) {
	// 1. Create webhook log entry
	webhookLog := &models.WebhookLog{
		ID:         primitive.NewObjectID(),
		EndpointID: endpointID,
		Method:     method,
		Headers:    headers,
		Body:       body,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Status:     "received",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 2. Store in database
	if err := s.StoreWebhookLog(webhookLog); err != nil {
		log.Printf("Failed to store webhook log: %v", err)
		webhookLog.Status = "failed"
		webhookLog.UpdatedAt = time.Now()
		return webhookLog, err
	}

	// 3. Trigger configured actions (async)
	go s.TriggerActions(webhookLog)

	// 4. Update status to processed
	webhookLog.Status = "processed"
	webhookLog.ProcessedAt = &time.Time{}
	*webhookLog.ProcessedAt = time.Now()
	webhookLog.UpdatedAt = time.Now()

	// Update the database record
	if err := s.UpdateWebhookLog(webhookLog); err != nil {
		log.Printf("Failed to update webhook log status: %v", err)
	}

	return webhookLog, nil
}

// StoreWebhookLog saves the webhook log to MongoDB
func (s *WebhookService) StoreWebhookLog(webhookLog *models.WebhookLog) error {
	collection := db.Mongo.Collection("webhook_logs")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, webhookLog)
	if err != nil {
		return err
	}

	log.Printf("Webhook log stored with ID: %s", webhookLog.ID.Hex())
	return nil
}

// UpdateWebhookLog updates an existing webhook log
func (s *WebhookService) UpdateWebhookLog(webhookLog *models.WebhookLog) error {
	collection := db.Mongo.Collection("webhook_logs")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": webhookLog.ID}
	update := bson.M{
		"$set": bson.M{
			"status":       webhookLog.Status,
			"processed_at": webhookLog.ProcessedAt,
			"updated_at":   webhookLog.UpdatedAt,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

// TriggerActions handles any configured actions for the webhook
func (s *WebhookService) TriggerActions(webhookLog *models.WebhookLog) {
	log.Printf("Triggering actions for webhook %s", webhookLog.ID.Hex())

	// TODO: Implement action triggering logic
	// This could include:
	// - Sending notifications
	// - Calling external APIs
	// - Triggering webhooks to other services
	// - Processing data transformations
	// - Logging to external systems

	// For now, just log the action
	log.Printf("Webhook %s processed successfully for endpoint %s",
		webhookLog.ID.Hex(), webhookLog.EndpointID)
}

// GetWebhookLogs retrieves webhook logs for an endpoint
func (s *WebhookService) GetWebhookLogs(endpointID string, limit int64) ([]models.WebhookLog, error) {
	collection := db.Mongo.Collection("webhook_logs")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"endpoint_id": endpointID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var webhookLogs []models.WebhookLog
	if err = cursor.All(ctx, &webhookLogs); err != nil {
		return nil, err
	}

	return webhookLogs, nil
}

// GetWebhookLogByID retrieves a specific webhook log by its ID
func (s *WebhookService) GetWebhookLogByID(webhookLogID string) (*models.WebhookLog, error) {
	collection := db.Mongo.Collection("webhook_logs")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(webhookLogID)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook log ID: %v", err)
	}

	filter := bson.M{"_id": objectID}

	var webhookLog models.WebhookLog
	err = collection.FindOne(ctx, filter).Decode(&webhookLog)
	if err != nil {
		return nil, fmt.Errorf("webhook log not found: %v", err)
	}

	return &webhookLog, nil
}

// ClearWebhookLogs deletes all webhook logs for an endpoint
func (s *WebhookService) ClearWebhookLogs(endpointID string) (int64, error) {
	collection := db.Mongo.Collection("webhook_logs")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"endpoint_id": endpointID}

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	log.Printf("Cleared %d webhook logs for endpoint: %s", result.DeletedCount, endpointID)
	return result.DeletedCount, nil
}
