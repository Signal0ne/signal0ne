package models

type IIntegration interface {
	Execute(input any,
		output map[string]string,
		functionName string) ([]map[string]any, error)

	Validate() error

	ValidateStep(input any,
		functionName string) error

	Initialize() map[string]string
}

type Integration struct {
	Id          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	NamespaceId string `json:"namespaceId" bson:"namespaceId"`
	Type        string `json:"type" bson:"type"`
}
