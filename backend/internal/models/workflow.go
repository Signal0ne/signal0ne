package models

type Workflow struct {
	NamespaceId  string  `bson:"namespaceId"`
	WorkflowId   string  `bson:"_id"`
	WorkflowSalt string  `bson:"salt"`
	Name         string  `json:"name" bson:"name"`
	Lookback     string  `json:"lookback" bson:"lookback"`
	Trigger      Trigger `json:"trigger" bson:"trigger"`
	Steps        []Step  `json:"steps" bson:"steps"`
}
