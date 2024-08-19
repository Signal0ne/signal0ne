package models

type EnrichedAlert struct {
	AdditionalContext map[string]Outputs `json:"additionalContext" bson:"additionalProperties"`
	TriggerProperties map[string]any     `json:"triggerProperties,inline" bson:"triggerProperties,inline"`
}

type Outputs struct {
	Output any `json:"output" bson:"output"`
}
