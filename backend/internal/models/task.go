package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Assignee User               `json:"assignee" bson:"assignee"`
	Comments []Comment          `json:"comments" bson:"comments"`
	IsDone   bool               `json:"isDone" bson:"isDone"`
	Items    []Item             `json:"items" bson:"items"`
	Priority int                `json:"priority" bson:"priority"`
	TaskName string             `json:"taskName" bson:"taskName"`
}

type Comment struct {
	Content   ItemContent `json:"content" bson:"content"`
	Source    User        `json:"source" bson:"source"`
	Timestamp int64       `json:"timestamp" bson:"timestamp"`
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
	Graph    ValueType = "graph"
	Markdown ValueType = "markdown"
)
