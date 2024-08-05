package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Workflow struct {
	Executions   []Execution        `json:"executions,omitempty" bson:"executions,omitempty"`
	Id           primitive.ObjectID `bson:"_id"`
	Lookback     string             `json:"lookback" bson:"lookback"`
	Name         string             `json:"name" bson:"name"`
	NamespaceId  string             `bson:"namespaceId"`
	Steps        []Step             `json:"steps" bson:"steps"`
	Trigger      Trigger            `json:"trigger" bson:"trigger"`
	WorkflowSalt string             `bson:"salt"`
}

type Execution struct {
	Log       string `json:"log" bson:"log"`
	Status    string `json:"status" bson:"status"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
}
