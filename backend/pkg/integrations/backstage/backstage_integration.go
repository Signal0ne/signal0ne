package backstage

import (
	"encoding/json"
	"fmt"
	"signal0ne/internal/models"
)

var functions = map[string]func(T any, dryRun bool) (any, error){
	"get_properties_values": getPropertiesValues,
}

type BackstageIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

func (integration BackstageIntegration) Execute(
	input interface{},
	output interface{},
	functionName string) (map[string]interface{}, error) {

	var result map[string]interface{}

	// [TBD]: execute function

	return result, nil
}

func (integration BackstageIntegration) Validate() error {
	if integration.Config.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if integration.Config.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	return nil
}

func (integration BackstageIntegration) ValidateStep(
	input interface{},
	output interface{},
	functionName string,
) error {
	function, exists := functions[functionName]
	if !exists {
		return fmt.Errorf("cannot find selected function")
	}

	_, err := function(input, true)
	if err != nil {
		return err
	}

	return nil
}

type GetPropertiesValuesInput struct {
	Filter string `json:"filter"`
}

func getPropertiesValues(T any, dryRun bool) (any, error) {
	var input GetPropertiesValuesInput
	data, err := json.Marshal(T)
	if err != nil {
		return nil, fmt.Errorf("invalid input for get_properties_values function")
	}

	err = json.Unmarshal(data, &input)
	if err != nil {
		return nil, fmt.Errorf("invalid input for get_properties_values function")
	}

	if dryRun {
		return nil, nil
	} else {
		// [TBD]: Execute
	}

	return input.Filter, nil
}
