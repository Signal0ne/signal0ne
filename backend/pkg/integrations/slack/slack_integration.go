package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
	"strings"
	"time"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"post_message": models.WorkflowFunctionDefinition{
		Function: postMessage,
		Input:    PostMessageInput{},
	},
}

type SlackIntegrationInventory struct {
	AlertTitle string
}

func NewSlackIntegrationInventory(alertTitle string) SlackIntegrationInventory {
	return SlackIntegrationInventory{
		AlertTitle: alertTitle,
	}
}

type SlackIntegration struct {
	Inventory          SlackIntegrationInventory
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

func (integration SlackIntegration) Execute(
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

func (integration SlackIntegration) Validate() error {
	if integration.Config.WorkspaceID == "" {
		return fmt.Errorf("host cannot be empty")
	}
	return nil
}

func (integration SlackIntegration) ValidateStep(
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

type PostMessageInput struct {
	SlackChannel          string `json:"slack_channel"`
	ParsableContextObject string `json:"parsable_context_object"`
	IgnoreContextKeys     string `json:"ignore_context_keys"`
}

func postMessage(input any, integration any) (output []any, err error) {
	var parsedInput PostMessageInput
	var parsedAlert models.EnrichedAlert

	err = helpers.ValidateInputParameters(input, &parsedInput, "post_message")
	if err != nil {
		return output, err
	}

	assertedIntegration := integration.(SlackIntegration)

	fmt.Printf("Executing slack postMessage\n")
	err = json.Unmarshal([]byte(parsedInput.ParsableContextObject), &parsedAlert)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	prettyJSON, err := json.MarshalIndent(parsedAlert, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Print(string(prettyJSON))

	url := fmt.Sprintf("http://%s:%s/api/post_message", assertedIntegration.Host, assertedIntegration.Port)
	title := assertedIntegration.Inventory.AlertTitle
	id := parsedAlert.Id.Hex()

	resultTime := time.Unix(parsedAlert.TriggerProperties["timestamp"].(int64), 0)
	date := resultTime.Format(time.RFC3339)

	data := map[string]any{
		"startedAt": date,
		"service":   parsedAlert.TriggerProperties["service"],
		"summary":   parsedAlert.AdditionalContext["openai_prod_summarize_context"].Output.([]any)[0].(map[string]any)["summary"],
	}

	payload := map[string]any{
		"channelName": strings.Split(parsedInput.SlackChannel, ",")[0],
		"title":       title,
		"id":          id,
		"data":        data,
	}

	prettyJSON, err = json.MarshalIndent(payload, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(prettyJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)

	}

	return output, err
}
