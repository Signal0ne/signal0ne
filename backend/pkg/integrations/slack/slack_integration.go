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
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"post_message": models.WorkflowFunctionDefinition{
		Function: postMessage,
		Input:    PostMessageInput{},
	},
	"create_channel": models.WorkflowFunctionDefinition{
		Function: createChannel,
		Input:    CreateChannelInput{},
	},
	"add_users_to_channel": models.WorkflowFunctionDefinition{
		Function: addUsersToTheChannel,
		Input:    AddUsersToTheChannelInput{},
	},
}

type SlackIntegrationInventory struct {
	AlertTitle string `json:"-" bson:"-"`
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
	IgnoreContextKeys     string `json:"ignore_context_keys"`
	ParsableContextObject string `json:"parsable_context_object"`
	PostMessagePayload    string `json:"post_message_payload"`
	SlackChannel          string `json:"slack_channel"`
}

type CreateChannelInput struct {
	ChannelName string `json:"channel_name"`
	IsPrivate   string `json:"is_private"`
}

type AddUsersToTheChannelInput struct {
	ChannelName string `json:"channel_name"`
	UserHandles string `json:"user_handles"`
}

func postMessage(input any, integration any) (output []any, err error) {
	var parsedInput PostMessageInput
	var parsedAlert models.EnrichedAlert

	err = helpers.ValidateInputParameters(input, &parsedInput, "post_message")
	if err != nil {
		return output, err
	}

	assertedIntegration := integration.(SlackIntegration)

	fmt.Printf("###\nExecuting slack postMessage\n")
	err = json.Unmarshal([]byte(parsedInput.ParsableContextObject), &parsedAlert)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	url := fmt.Sprintf("http://%s:%s/api/post_message", assertedIntegration.Host, assertedIntegration.Port)
	title := assertedIntegration.Inventory.AlertTitle
	id := parsedAlert.Id.Hex()

	data := map[string]any{}
	err = json.Unmarshal([]byte(parsedInput.PostMessagePayload), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	payload := map[string]any{
		"channelName": strings.Split(parsedInput.SlackChannel, ",")[0],
		"data":        data,
		"id":          id,
		"title":       title,
	}

	prettyJSON, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(prettyJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)

	}

	return output, err
}

func createChannel(input any, integration any) (output []any, err error) {

	var parsedInput CreateChannelInput

	err = helpers.ValidateInputParameters(input, &parsedInput, "create_channel")
	if err != nil {
		return output, err
	}

	assertedIntegration := integration.(SlackIntegration)
	url := fmt.Sprintf("http://%s:%s/api/create_channel", assertedIntegration.Host, assertedIntegration.Port)

	payload := map[string]any{
		"channelName": parsedInput.ChannelName,
		"isPrivate":   parsedInput.IsPrivate,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create channel: %v", resp.Status)
	}

	var response map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	output = append(output, response)
	return output, err
}

func addUsersToTheChannel(input any, integration any) (output []any, err error) {

	var parsedInput AddUsersToTheChannelInput

	err = helpers.ValidateInputParameters(input, &parsedInput, "add_users_to_channel")
	if err != nil {
		return output, err
	}

	assertedIntegration := integration.(SlackIntegration)
	url := fmt.Sprintf("http://%s:%s/api/add_users_to_channel", assertedIntegration.Host, assertedIntegration.Port)

	payload := map[string]any{
		"channelName": parsedInput.ChannelName,
		"userHandles": strings.Split(parsedInput.UserHandles, ","),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to add users to the channel: %v", resp.Status)
	}

	var response map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	output = append(output, response)
	return output, err
}
