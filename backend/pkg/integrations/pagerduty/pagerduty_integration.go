package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
)

type PagerdutyIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

var functions = map[string]models.WorkflowFunctionDefinition{
	"create_incident": models.WorkflowFunctionDefinition{
		Function: createIncident,
		Input:    CreateIncidentInput{},
	},
}

func (integration PagerdutyIntegration) Execute(
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

func (integration PagerdutyIntegration) Validate() error {
	// Implement your config validation here
	return nil
}

func (integration PagerdutyIntegration) ValidateStep(
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

type CreateIncidentInput struct {
	Type        string `json:"type" bson:"type"`
	Title       string `json:"title" bson:"title"`
	Urgency     string `json:"urgency" bson:"urgency"`
	ServiceName string `json:"service_name" bson:"service_name"`
}

func createIncident(input any, integration any) ([]any, error) {
	var parsedInput CreateIncidentInput
	var output []any

	assertedIntegration := integration.(PagerdutyIntegration)

	err := helpers.ValidateInputParameters(input, &parsedInput, "create_incident")
	if err != nil {
		return output, err
	}

	client := http.Client{}

	//Search for the service by name

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%/services?service=%s", assertedIntegration.Config.Url, parsedInput.ServiceName),
		nil,
	)

	req.Header.Add("Authorization", fmt.Sprintf("Token token=%s", assertedIntegration.Config.ApiKey))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	var services map[string]any
	servicesResponse, err := client.Do(req)
	if err != nil {
		return output, err
	}

	err = json.NewDecoder(servicesResponse.Body).Decode(&services)
	if err != nil {
		return output, err
	}
	defer servicesResponse.Body.Close()

	var service = make(map[string]string)

	service["id"] = services["services"].([]any)[0].(map[string]any)["id"].(string)
	service["type"] = services["services"].([]any)[0].(map[string]any)["type"].(string)

	//Create incident
	var incidentBody = map[string]any{
		"type":    parsedInput.Type,
		"title":   parsedInput.Title,
		"urgency": parsedInput.Urgency,
		"service": service,
	}

	incidentBodyJSON, err := json.Marshal(incidentBody)
	if err != nil {
		return output, err
	}

	incidentBodyBuffer := bytes.NewBuffer(incidentBodyJSON)

	req, err = http.NewRequest("POST", fmt.Sprintf("%s/incidents", assertedIntegration.Config.Url), incidentBodyBuffer)
	if err != nil {
		return output, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token token=%s", assertedIntegration.Config.ApiKey))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	incidentResponse, err := client.Do(req)
	if err != nil {
		return output, nil
	}

	output = append(output, incidentResponse)

	return output, nil
}
