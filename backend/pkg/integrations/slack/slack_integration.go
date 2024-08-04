package slack

import (
	"encoding/json"
	"fmt"
	"signal0ne/internal/models"
	"signal0ne/pkg/integrations/helpers"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"post_message": models.WorkflowFunctionDefinition{
		Function: postMessage,
		Input:    PostMessageInput{},
	},
}

type SlackIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

func (integration SlackIntegration) Execute(
	input any,
	output map[string]string,
	functionName string) ([]any, error) {

	var result []any

	function, ok := functions[functionName]
	if !ok {
		return result, fmt.Errorf("cannot find requested function")
	}

	result, err := function.Function(input)
	if err != nil {
		return make([]any, 0), err
	}

	return result, nil
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
	ParsableContextObject string `json:"parsable_context_object"`
	IgnoreContextKeys     string `json:"ignore_context_keys"`
}

func postMessage(input any) (output []any, err error) {
	var parsedInput PostMessageInput
	var parsedAlert map[string]any

	err = helpers.ValidateInputParameters(input, &parsedInput, "post_message")
	if err != nil {
		return output, err
	}

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

	return output, err
}
