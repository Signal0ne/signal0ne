package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type IncidentStatus string

const (
	IncidentStatusOpen     IncidentStatus = "open"
	IncidentStatusResolved IncidentStatus = "resolved"
)

type IncidentSeverity string

const (
	IncidentSeverityCritical IncidentSeverity = "critical"
	IncidentSeverityHigh     IncidentSeverity = "high"
	IncidentSeverityModerate IncidentSeverity = "moderate"
	IncidentSeverityLow      IncidentSeverity = "low"
)

type Incident struct {
	Id          primitive.ObjectID       `json:"id" bson:"_id"`
	Assignee    User                     `json:"assignee" bson:"assignee"`
	History     []IncidentUpdate[Update] `json:"history" bson:"history"`
	NamespaceId string                   `json:"namespaceId" bson:"namespaceId"`
	Severity    IncidentSeverity         `json:"severity" bson:"severity"`
	Status      IncidentStatus           `json:"status" bson:"status"`
	Summary     string                   `json:"summary" bson:"summary"`
	Tasks       []Task                   `json:"tasks" bson:"tasks"`
	Timestamp   int64                    `json:"timestamp" bson:"timestamp"`
	Title       string                   `json:"title" bson:"title"`
}

type IncidentUpdate[T Update] struct {
	Doer      User               `json:"doer" bson:"doer"`
	Timestamp primitive.DateTime `json:"timestamp" bson:"timestamp"`
	Type      string             `json:"type" bson:"type"`
	Update    T                  `json:"update" bson:"update"`
}

type Update interface {
	IsUpdate()
}

type AssigneeUpdate struct {
	New User `json:"new" bson:"new"`
	Old User `json:"old" bson:"old"`
}

func (AssigneeUpdate) IsUpdate() {}

type TaskUpdate struct {
	FieldKey     string `json:"fieldKey" bson:"fieldKey"`
	TaskStepName string `json:"stepName" bson:"stepName"`
}

func (TaskUpdate) IsUpdate() {}
