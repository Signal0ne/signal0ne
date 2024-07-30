package models

type IIntegration interface {
	Execute() map[string]interface{}
	Validate() error
	ValidateStep() error
}

type Integration struct {
	Name     string                 `json:"name" bson:"name"`
	Type     string                 `json:"type" bson:"type"`
	ImageURL string                 `json:"imageUrl" bson:"imageUrl"`
	config   map[string]interface{} `json:"-"`
}

func (i *Integration) GetConfig() map[string]interface{} {
	return i.config
}
