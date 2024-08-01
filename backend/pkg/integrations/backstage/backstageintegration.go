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

func (i BackstageIntegration) Execute(
	input interface{},
	output interface{},
	functionName string) (map[string]interface{}, error) {

	var result map[string]interface{}

	// [TBD]: execute funtion

	return result, nil
}

func (i BackstageIntegration) Validate() error {
	if i.Config.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if i.Config.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	return nil
}

func (i BackstageIntegration) ValidateStep(
	input interface{},
	output interface{},
	functionName string,
) error {
	function, exists := functions[functionName]
	if !exists {
		return fmt.Errorf("cannot find selected funtion")
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
