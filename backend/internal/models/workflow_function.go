package models

type WorkflowFunctionDefinition struct {
	Function func(input any) (output []any, err error)
	Input    any
}
