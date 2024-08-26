package models

type Task struct {
	StepName string
	Priority int
	Assignee User
	IsDone   bool
	Items    []Item
}

type Item struct {
	Fields []Field
}

type Field struct {
	Key       string
	Source    string
	Value     any
	ValueType string
}
