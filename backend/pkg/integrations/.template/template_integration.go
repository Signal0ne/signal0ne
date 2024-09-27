//go:build ignore
// +build ignore

package template

import (
	"fmt"
	"signal0ne/internal/models"
	"signal0ne/pkg/integrations/helpers"
)

type TemplateIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

var functions = map[string]models.WorkflowFunctionDefinition{
	//Add functions provided by the integration
}

func (integration TemplateIntegration) Execute(
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

func (integration TemplateIntegration) Initialize() map[string]string {
	// Implement your config initialization here
	return nil
}

func (integration TemplateIntegration) Validate() error {
	// Implement your config validation here
	return nil
}

func (integration TemplateIntegration) ValidateStep(
	input any,
	functionName string,
) error {
	function, exists := functions[functionName]
	if !exists {
		return fmt.Errorf("cannot find selected function")
	}

	//Validate input parameters for the chosen function
	err := helpers.ValidateInputParameters(input, function.Input, functionName)
	if err != nil {
		return err
	}

	return nil
}

//Implement functions and it's input types below

//----------------------------------------------
