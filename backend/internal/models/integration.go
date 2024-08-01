package models

type IIntegration interface {
	Execute(input interface{},
		output interface{},
		functionName string) (map[string]interface{}, error)

	Validate() error

	ValidateStep(input interface{},
		output interface{},
		functionName string) error
}

type Integration struct {
	Name     string `json:"name" bson:"name"`
	Type     string `json:"type" bson:"type"`
	ImageURL string `json:"imageUrl" bson:"imageUrl"`
}
