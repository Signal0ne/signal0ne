package models

type Task struct {
	StepName string
	Priority int
	Assignee User
	IsDone   bool
	Fields   []Field
}

type Field struct {
	Key       string
	Source    string
	Value     any
	ValueType string
}
