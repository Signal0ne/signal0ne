package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlertStatus string

const (
	AlertStatusActive   AlertStatus = "open"
	AlertStatusInactive AlertStatus = "inactive"
)

type EnrichedAlert struct {
	Id                primitive.ObjectID `json:"id" bson:"_id"`
	AdditionalContext map[string]any     `json:"additionalContext" bson:"additionalProperties"`
	AlertName         string             `json:"alertName" bson:"alertName"`
	Integration       string             `json:"integration" bson:"integration"`
	OriginalUrl       string             `json:"originalUrl" bson:"originalUrl"`
	Tags              []string           `json:"tags" bson:"tags"`
	StartTime         time.Time          `json:"startTime" bson:"startTime"`
	State             AlertStatus        `json:"state" bson:"state"`
	Service           string             `json:"service" bson:"service"`
	TriggerProperties map[string]any     `json:"triggerProperties,inline" bson:"triggerProperties,inline"`
	WorkflowId        string             `json:"workflowId" bson:"workflowId"`
}
