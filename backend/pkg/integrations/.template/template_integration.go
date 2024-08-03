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
	functionName string) ([]any, error) {

	var result []any

	function, ok := functions[functionName]
	if !ok {
		return result, fmt.Errorf("cannot find requested function")
	}

	// Call proper function based on user workflow definition
	result, err := function.Function(input)
	if err != nil {
		return make([]any, 0), err
	}

	return result, nil
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

	//Valiudate input parameters for the chosen fucntion
	err := helpers.ValidateInputParameters(input, function.Input)
	if err != nil {
		return err
	}

	return nil
}

//Implement functions and it's input types below

//----------------------------------------------
