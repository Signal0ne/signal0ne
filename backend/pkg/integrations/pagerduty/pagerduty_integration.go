package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"create_incident": models.WorkflowFunctionDefinition{
		Function: createIncident,
		Input:    CreateIncidentInput{},
	},
}

type PagerdutyIntegrationInventory struct {
	WorkflowProperties *models.Workflow `json:"-" bson:"-"`
}

func NewPagerdutyIntegrationInventory(
	workflowProperties *models.Workflow) PagerdutyIntegrationInventory {
	return PagerdutyIntegrationInventory{
		WorkflowProperties: workflowProperties,
	}
}

type PagerdutyIntegration struct {
	Config             `json:",inline" bson:",inline"`
	Inventory          PagerdutyIntegrationInventory
	models.Integration `json:",inline" bson:",inline"`
}

type Note struct {
	Content string `json:"content"`
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
	Type                  string `json:"type" bson:"type"`
	Title                 string `json:"title" bson:"title"`
	Urgency               string `json:"urgency" bson:"urgency"`
	ServiceName           string `json:"service_name" bson:"service_name"`
	ParsableContextObject string `json:"parsable_context_object"`
}

func createIncident(input any, integration any) ([]any, error) {
	var parsedInput CreateIncidentInput
	var parsedAlert models.EnrichedAlert
	var output []any

	assertedIntegration := integration.(PagerdutyIntegration)

	err := helpers.ValidateInputParameters(input, &parsedInput, "create_incident")
	if err != nil {
		return output, err
	}

	err = json.Unmarshal([]byte(parsedInput.ParsableContextObject), &parsedAlert)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	client := http.Client{}

	// Search for the service by name

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/services?name=%s", assertedIntegration.Config.Url, url.QueryEscape(parsedInput.ServiceName)),
		nil,
	)
	if err != nil {
		return output, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token token=%s", assertedIntegration.Config.ApiKey))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	var services map[string]any
	servicesResponse, err := client.Do(req)
	if err != nil {
		return output, err
	}

	if servicesResponse.StatusCode != http.StatusOK {
		return output, fmt.Errorf("failed to get services: %v", servicesResponse.Status)
	}

	err = json.NewDecoder(servicesResponse.Body).Decode(&services)
	if err != nil {
		fmt.Printf("Error decoding response: %v", err)
		return output, err
	}
	defer servicesResponse.Body.Close()

	var service = make(map[string]string)

	service["id"] = services["services"].([]any)[0].(map[string]any)["id"].(string)
	service["type"] = services["services"].([]any)[0].(map[string]any)["type"].(string)

	// Create incident
	var incidentBody = map[string]any{
		"incident": map[string]any{
			"service": service,
			"title":   parsedInput.Title,
			"type":    parsedInput.Type,
		},
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
	req.Header.Add("From", "contact@signaloneai.com")

	incidentResponse, err := client.Do(req)
	if err != nil {
		return output, nil
	}

	if incidentResponse.StatusCode != http.StatusCreated {
		var errorMessage any
		err := json.NewDecoder(incidentResponse.Body).Decode(&errorMessage)
		if err != nil {
			return output, err
		}
		return output, fmt.Errorf("failed to create incident: %v", errorMessage)
	}

	incident := make(map[string]any)
	err = json.NewDecoder(incidentResponse.Body).Decode(&incident)
	if err != nil {
		return output, nil
	}

	incidentId, _ := incident["incident"].(map[string]any)["id"].(string)

	// Create note for each step from the workflow for the current incident
	for _, step := range assertedIntegration.Inventory.WorkflowProperties.Steps {
		var note = Note{}
		isDone := true

		stepOutputs, _ := parsedAlert.AdditionalContext[fmt.Sprintf("%s_%s", step.Integration, step.Function)].Output.([]any)

		if len(stepOutputs) == 0 {
			isDone = false
		}

		if !isDone {
			continue
		}

		note.Content = fmt.Sprintf("Step: %s\nAssignee: Signal0ne\n##############################\n", step.Name)

		for _, stepOutput := range stepOutputs {
			for key, value := range stepOutput.(map[string]any) {
				note.Content += fmt.Sprintf("Key: %s\nValue: %v\n---\n", key, value)
			}
		}

		note.Content += "##############################\n"

		noteBodyJSON, err := json.Marshal(map[string]Note{
			"note": note,
		})
		if err != nil {
			return output, err
		}

		noteBodyBuffer := bytes.NewBuffer(noteBodyJSON)

		req, err = http.NewRequest("POST", fmt.Sprintf("%s/incidents/%s/notes", assertedIntegration.Config.Url, incidentId), noteBodyBuffer)
		if err != nil {
			return output, err
		}

		req.Header.Add("Authorization", fmt.Sprintf("Token token=%s", assertedIntegration.Config.ApiKey))
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("From", "contact@signaloneai.com")

		_, err = client.Do(req)
		if err != nil {
			return output, err
		}
	}

	return output, nil
}
