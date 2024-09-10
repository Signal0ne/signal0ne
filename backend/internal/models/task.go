package models

type Task struct {
	TaskName string `json:"taskName" bson:"taskName"`
	Priority int    `json:"priority" bson:"priority"`
	Assignee User   `json:"assignee" bson:"assignee"`
	IsDone   bool   `json:"isDone" bson:"isDone"`
	Items    []Item `json:"items" bson:"items"`
}

type Item struct {
	Content []ItemContent `json:"content" bson:"content"`
	Source  string        `json:"source" bson:"source"`
}

type ItemContent struct {
	Key       string    `json:"key" bson:"key"`
	Value     any       `json:"value" bson:"value"`
	ValueType ValueType `json:"valueType" bson:"valueType"`
}

type ValueType string

const (
	Markdown ValueType = "markdown"
	Graph    ValueType = "graph"
)
