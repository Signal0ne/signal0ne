package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Incident struct {
	Id       primitive.ObjectID       `json:"id" bson:"_id"`
	Title    string                   `json:"title" bson:"title"`
	Assignee User                     `json:"assignee" bson:"assignee"`
	Severity string                   `json:"severity" bson:"severity"`
	Summary  string                   `json:"summary" bson:"summary"`
	Tasks    []Task                   `json:"tasks" bson:"tasks"`
	History  []IncidentUpdate[Update] `json:"history" bson:"history"`
}

type IncidentUpdate[T Update] struct {
	Timestamp primitive.DateTime `json:"timestamp" bson:"timestamp"`
	Type      string             `json:"type" bson:"type"`
	Doer      User               `json:"doer" bson:"doer"`
	Update    T                  `json:"update" bson:"update"`
}

type Update interface {
	IsUpdate()
}

type AssigneeUpdate struct {
	Old User `json:"old" bson:"old"`
	New User `json:"new" bson:"new"`
}

func (AssigneeUpdate) IsUpdate() {}

type TaskUpdate struct {
	TaskStepName string `json:"stepName" bson:"stepName"`
	FieldKey     string `json:"fieldKey" bson:"fieldKey"`
}

func (TaskUpdate) IsUpdate() {}
