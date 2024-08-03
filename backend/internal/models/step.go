package models

type Step struct {
	Name        string            `json:"name" bson:"name"`
	Function    string            `json:"function" bson:"function"`
	Integration string            `json:"integration" bson:"integration"`
	Input       map[string]string `json:"input" bson:"input"`
	Output      map[string]string `json:"output" bson:"output"`
	Condition   string            `json:"condition" bson:"condition"`
}
