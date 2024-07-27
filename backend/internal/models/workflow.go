package models

type Workflow struct {
	NamespaceId string `bson:"namespaceId"`
	WorkflowId  string `bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Lookback    string `json:"lookback" bson:"lookback"`
	Steps       []Step `json:"steps" bson:"steps"`
}
