package models

type Integration struct {
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	ImageURL string                 `json:"imageUrl"`
	config   map[string]interface{} `json:"-"`
}

func (i *Integration) GetConfig() map[string]interface{} {
	return i.config
}
