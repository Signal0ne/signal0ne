package jaeger

import (
	"fmt"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"get_properties_values": models.WorkflowFunctionDefinition{
		Function: getPropertiesValues,
		Input:    GetPropertiesValuesInput{},
	},
}

type JaegerIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

func (integration JaegerIntegration) Execute(
	input any,
	output map[string]string,
	functionName string) ([]map[string]any, error) {

	var results []map[string]any

	function, ok := functions[functionName]
	if !ok {
		return results, fmt.Errorf("%s.%s: cannot find requested function", integration.Name, functionName)
	}

	intermediateResults, err := function.Function(input, integration)
	if err != nil {
		return results, fmt.Errorf("%s.%s:%v", integration.Name, functionName, err)
	}

	results = tools.ExecutionResultWrapper(intermediateResults, output)

	return results, nil
}

func (integration JaegerIntegration) Validate() error {
	if integration.Config.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if integration.Config.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	return nil
}

func (integration JaegerIntegration) ValidateStep(
	input any,
	functionName string,
) error {
	function, exists := functions[functionName]
	if !exists {
		return fmt.Errorf("cannot find selected function")
	}

	err := helpers.ValidateInputParameters(input, function.Input, functionName)
	if err != nil {
		return err
	}

	return nil
}

type GetPropertiesValuesInput struct {
	Service string `json:"service"`
	Tags    string `json:"tags"`
}

func getPropertiesValues(input any, integration any) ([]any, error) {
	var parsedInput GetPropertiesValuesInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "get_properties_values")
	if err != nil {
		return output, err
	}

	fmt.Printf("Executing Jaeger integration function...")

	return output, nil

}
