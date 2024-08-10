package signal0ne

import (
	"fmt"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
)

type Signal0neIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

var functions = map[string]models.WorkflowFunctionDefinition{
	"correlate_ongoing_alerts": models.WorkflowFunctionDefinition{
		Function: correlateOngoingAlerts,
		Input:    CorrelateOngoingAlertsInput{},
	},
}

func (integration Signal0neIntegration) Execute(
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

func (integration Signal0neIntegration) Validate() error {
	// Implement your config validation here
	return nil
}

func (integration Signal0neIntegration) ValidateStep(
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

type CorrelateOngoingAlertsInput struct {
	Filter string `json:"filter"`
}

func correlateOngoingAlerts(input any, integration any) ([]any, error) {
	var parsedInput CorrelateOngoingAlertsInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "correlate_ongoing_alerts")
	if err != nil {
		return output, err
	}

	fmt.Printf("Executing backstage getPropertiesValues\n")
	// 1.Get alerts by filter
	// 2.Run semantic similarity
	// 3.Return similar alert objects
	// What do we do with resolved alerts

	return output, nil
}
