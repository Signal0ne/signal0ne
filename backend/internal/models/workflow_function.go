package models

type WorkflowFunctionDefinition struct {
	Function func(input any, integration any) (output []any, err error)
	Input    any
}
