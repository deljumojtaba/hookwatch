package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebhookLog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	EndpointID  string             `bson:"endpoint_id" json:"endpoint_id"`
	Method      string             `bson:"method" json:"method"`
	Headers     map[string]string  `bson:"headers" json:"headers"`
	Body        interface{}        `bson:"body" json:"body"`
	IPAddress   string             `bson:"ip_address" json:"ip_address"`
	UserAgent   string             `bson:"user_agent" json:"user_agent"`
	Status      string             `bson:"status" json:"status"` // "received", "processed", "failed"
	ProcessedAt *time.Time         `bson:"processed_at,omitempty" json:"processed_at,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type WebhookEndpoint struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Secret      string             `bson:"secret" json:"secret,omitempty"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
