package backstage

import (
	"encoding/json"
	"signal0ne/internal/models"
)

type BackstageIntegration struct {
	models.Integration
	Config Config
}

func NewBackstageIntegration(integrationTemplate models.Integration) BackstageIntegration {
	var config Config
	jsonString, err := json.Marshal(integrationTemplate.GetConfig())
	if err != nil {
		return BackstageIntegration{}
	}
	json.Unmarshal(jsonString, &config)
	return BackstageIntegration{
		Integration: integrationTemplate,
		Config:      config,
	}
}

func (i *BackstageIntegration) Execute(
	input map[string]string,
	output map[string]string,
	mapping map[string]string) map[string]string {

	return make(map[string]string)
}

func (i *BackstageIntegration) Validate() bool {
	return false
}
