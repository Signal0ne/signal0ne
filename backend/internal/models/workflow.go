package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Workflow struct {
	Description  string             `json:"description" bson:"description"`
	Executions   []StepExecution    `json:"executions,omitempty" bson:"executions,omitempty"`
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	Lookback     string             `json:"lookback" bson:"lookback"`
	Name         string             `json:"name" bson:"name"`
	NamespaceId  string             `json:"namespaceId" bson:"namespaceId"`
	Steps        []Step             `json:"steps" bson:"steps"`
	Trigger      Trigger            `json:"trigger" bson:"trigger"`
	WorkflowSalt string             `json:"salt" bson:"salt"`
}

type StepExecution struct {
	ParsedWorkflow ParsedWorkflow         `json:"parsedWorkflow" bson:"parsedWorkflow"`
	Outputs        map[string]any         `json:"outputs" bson:"outputs"`
	Outcomes       []StepExecutionOutcome `json:"outcomes" bson:"outcomes"`
}

type StepExecutionOutcome struct {
	Status     string `json:"status" bson:"status"`
	LogMessage string `json:"logMessage" bson:"logMessage"`
}

type ParsedWorkflow struct {
	Steps   []Step  `json:"steps" bson:"steps"`
	Trigger Trigger `json:"trigger" bson:"trigger"`
}
