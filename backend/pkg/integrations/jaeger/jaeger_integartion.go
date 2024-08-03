package jaeger

import (
	"fmt"
	"signal0ne/internal/models"
	"signal0ne/pkg/integrations/helpers"
)

var functions = map[string]models.WorkflowFunctionDefinition{}

type JaegerIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

func (integration JaegerIntegration) Execute(
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
