package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AlertStatus string

const (
	AlertStatusActive   AlertStatus = "active"
	AlertStatusInactive AlertStatus = "inactive"
)

type EnrichedAlert struct {
	Id                primitive.ObjectID `json:"id" bson:"_id"`
	AdditionalContext map[string]any     `json:"additionalContext" bson:"additionalProperties"`
	AlertName         string             `json:"alertName" bson:"alertName"`
	Integration       string             `json:"integration" bson:"integration"`
	StartTime         string             `json:"startTime" bson:"startTime"`
	State             AlertStatus        `json:"state" bson:"state"`
	TriggerProperties map[string]any     `json:"triggerProperties,inline" bson:"triggerProperties,inline"`
	WorkflowId        string             `json:"workflowId" bson:"workflowId"`
}
