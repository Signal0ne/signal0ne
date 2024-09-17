package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AlertStatus string

const (
	AlertStatusActive   AlertStatus = "active"
	AlertStatusInactive AlertStatus = "inactive"
)

type EnrichedAlert struct {
	Id                primitive.ObjectID `json:"id" bson:"_id"`
	WorkflowId        string             `json:"workflowId" bson:"workflowId"`
	State             AlertStatus        `json:"state" bson:"state"`
	Integration       string             `json:"integration" bson:"integration"`
	AdditionalContext map[string]any     `json:"additionalContext" bson:"additionalProperties"`
	TriggerProperties map[string]any     `json:"triggerProperties,inline" bson:"triggerProperties,inline"`
}
