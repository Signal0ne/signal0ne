package confluence

import (
	"fmt"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
)

type ConfluenceIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

var functions = map[string]models.WorkflowFunctionDefinition{
	"search": models.WorkflowFunctionDefinition{
		Function: search,
		Input:    SearchInput{},
	},
}

func (integration ConfluenceIntegration) Execute(
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

func (integration ConfluenceIntegration) Validate() error {
	// Implement your config validation here
	return nil
}

func (integration ConfluenceIntegration) ValidateStep(
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

type SearchInput struct {
	Query string `json:"query" bson:"query"`
}

func search(input any, integration any) ([]any, error) {
	var parsedInput SearchInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "get_content")
	if err != nil {
		return output, err
	}

	return output, nil

}
