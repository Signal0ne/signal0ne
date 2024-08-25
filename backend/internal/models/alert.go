package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EnrichedAlert struct {
	Id                primitive.ObjectID `json:"id" bson:"_id"`
	AdditionalContext map[string]Outputs `json:"additionalContext" bson:"additionalProperties"`
	TriggerProperties map[string]any     `json:"triggerProperties,inline" bson:"triggerProperties,inline"`
}

type Outputs struct {
	Output any `json:"output" bson:"output"`
}
