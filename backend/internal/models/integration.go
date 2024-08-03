package models

type IIntegration interface {
	Execute(input any,
		output map[string]string,
		functionName string) ([]any, error)

	Validate() error

	ValidateStep(input any,
		functionName string) error
}

type Integration struct {
	Name string `json:"name" bson:"name"`
	Type string `json:"type" bson:"type"`
}
