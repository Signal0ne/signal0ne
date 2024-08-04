package models

type EnrichedAlert struct {
	TriggerProperties map[string]any     `json:"triggerProperties,inline" bson:"triggerProperties,inline"`
	AdditionalContext map[string]Outputs `json:"additionalContext" bson:"additionalProperties"`
}

type Outputs struct {
	Output []any `json:"output" bson:"output"`
}
