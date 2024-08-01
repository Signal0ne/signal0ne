package slack

import (
	"encoding/json"
	"fmt"
	"signal0ne/internal/models"
)

var functions = map[string]func(T any, dryRun bool) (any, error){
	"post_message": postMessage,
}

type SlackIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

func (i SlackIntegration) Execute(
	input interface{},
	output interface{},
	functionName string) (map[string]interface{}, error) {

	var result map[string]interface{}

	// [TBD]: execute funtion

	return result, nil
}

func (i SlackIntegration) Validate() error {
	if i.Config.WorkspaceID == "" {
		return fmt.Errorf("host cannot be empty")
	}
	return nil
}

func (i SlackIntegration) ValidateStep(
	input interface{},
	output interface{},
	functionName string,
) error {
	function, exists := functions[functionName]
	if !exists {
		return fmt.Errorf("cannot find selected function")
	}

	_, err := function(input, true)
	if err != nil {
		return err
	}

	return nil
}

type PostMessageInput struct {
	ParsableContextObject string   `json:"parsable_context_object"`
	IgnoreContextKeys     []string `json:"ignore_context_keys"`
}

func postMessage(T any, dryRun bool) (any, error) {
	var input PostMessageInput
	data, err := json.Marshal(T)
	if err != nil {
		return nil, fmt.Errorf("invalid input for post_message function")
	}

	err = json.Unmarshal(data, &input)
	if err != nil {
		return nil, fmt.Errorf("invalid input for post_message function")
	}

	if dryRun {
		return nil, nil
	} else {
		// [TBD]: Execute
	}

	return nil, nil
}
