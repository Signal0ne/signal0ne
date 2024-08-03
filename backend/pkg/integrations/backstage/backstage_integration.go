package backstage

import (
	"fmt"
	"signal0ne/internal/models"
	"signal0ne/pkg/integrations/helpers"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"get_properties_values": models.WorkflowFunctionDefinition{
		Function: getPropertiesValues,
		Input:    GetPropertiesValuesInput{},
	},
}

type BackstageIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

func (integration BackstageIntegration) Execute(
	input any,
	output map[string]string,
	functionName string) ([]any, error) {

	var result []any

	function, ok := functions[functionName]
	if !ok {
		return result, fmt.Errorf("cannot find requested function")
	}

	result, err := function.Function(input)
	if err != nil {
		return make([]any, 0), err
	}

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
	input any,
	functionName string,
) error {
	function, exists := functions[functionName]
	if !exists {
		return fmt.Errorf("cannot find selected function")
	}

	err := helpers.ValidateInputParameters(input, function.Input)
	if err != nil {
		return err
	}

	return nil
}

type GetPropertiesValuesInput struct {
	Filter string `json:"filter"`
}

func getPropertiesValues(input any) (output []any, err error) {
	var parsedInput GetPropertiesValuesInput

	err = helpers.ValidateInputParameters(input, parsedInput)
	if err != nil {
		return output, err
	}

	return output, err
}
