package models

type Step struct {
	Condition   string            `json:"condition,omitempty" bson:"condition"`
	Function    string            `json:"function" bson:"function"`
	Input       map[string]string `json:"input" bson:"input"`
	Integration string            `json:"integration" bson:"integration"`
	Name        string            `json:"name" bson:"name"`
	Output      map[string]string `json:"output,omitempty" bson:"output"`
}
