package models

type Step struct {
	Name        string      `json:"name" bson:"name"`
	Action      string      `json:"action" bson:"action"`
	Integration string      `json:"integration" bson:"integration"`
	Config      interface{} `json:"config" bson:"config"`
}
