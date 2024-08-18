package alertmanager

import (
	"context"
	"fmt"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"

	"github.com/prometheus/alertmanager/api/v2/client"
	"github.com/prometheus/alertmanager/api/v2/client/alert"
	alertmanagerApiModels "github.com/prometheus/alertmanager/api/v2/models"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"get_relevant_alerts": models.WorkflowFunctionDefinition{
		Function: getRelevantAlerts,
		Input:    GetRelevantAlertsInput{},
	},
}

type AlertmanagerIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

func (integration AlertmanagerIntegration) Execute(
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

func (integration AlertmanagerIntegration) Validate() error {
	if integration.Config.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if integration.Config.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	return nil
}

func (integration AlertmanagerIntegration) ValidateStep(
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

type GetRelevantAlertsInput struct {
	Filter string `json:"filter"`
}

func getRelevantAlerts(input any, integration any) ([]any, error) {
	var parsedInput GetRelevantAlertsInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "compare_traces")
	if err != nil {
		return output, err
	}

	fmt.Printf("Executing Jaeger integration function...")

	assertedIntegration := integration.(AlertmanagerIntegration)

	host := assertedIntegration.Host
	port := assertedIntegration.Port

	transport := client.DefaultTransportConfig().
		WithHost(fmt.Sprintf("%s:%s", host, port)).
		WithSchemes([]string{"http"})

	_ = client.NewHTTPClientWithConfig(nil, transport)

	fmt.Printf("PARSED ALERT FILTERS: %v", parsedInput.Filter)

	// alerts, err := getAlerts(alertmanagerClient, parsedInput.Filter)

	return output, nil
}

func getAlerts(c *client.AlertmanagerAPI, filters []string) ([]*alertmanagerApiModels.GettableAlert, error) {
	params := alert.NewGetAlertsParams().WithContext(context.Background())

	if len(filters) > 0 {
		params.SetFilter(filters)
	}

	result, err := c.Alert.GetAlerts(params)
	if err != nil {
		return nil, err
	}

	return result.Payload, nil
}
