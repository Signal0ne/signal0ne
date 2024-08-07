package backstage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type BackstageIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

func (integration BackstageIntegration) Execute(
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

	err := helpers.ValidateInputParameters(input, function.Input, functionName)
	if err != nil {
		return err
	}

	return nil
}

type GetPropertiesValuesInput struct {
	Filter string `json:"filter"`
}

func getPropertiesValues(input any, integration any) ([]any, error) {
	var parsedInput GetPropertiesValuesInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "get_properties_values")
	if err != nil {
		return output, err
	}

	fmt.Printf("Executing backstage getPropertiesValues\n")

	assertedIntegration := integration.(BackstageIntegration)

	host := assertedIntegration.Config.Host
	port := assertedIntegration.Config.Port
	apiKey := assertedIntegration.Config.ApiKey
	apiPath := "/api/catalog/entities/by-query?filter="

	url := fmt.Sprintf("http://%s:%s%s%s", host, port, apiPath, parsedInput.Filter)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return output, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return output, err
	}
	defer resp.Body.Close()
	var bodyHandler map[string]any

	if resp.StatusCode != 200 {
		err = fmt.Errorf("%s", resp.Status)
		return output, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return output, err
	}

	json.Unmarshal(body, &bodyHandler)
	items, exists := bodyHandler["items"]
	if !exists {
		err = fmt.Errorf("cannot parse %s response body", assertedIntegration.Name)
		return output, err
	}

	output = items.([]any)

	return output, err
}
