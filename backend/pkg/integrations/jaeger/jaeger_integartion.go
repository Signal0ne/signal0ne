package jaeger

import (
	"fmt"
	"signal0ne/internal/models"
)

var functions = map[string]func(T any, dryRun bool) (any, error){}

type JaegerIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

func (integration JaegerIntegration) Execute(
	input interface{},
	output interface{},
	functionName string) (map[string]interface{}, error) {

	var result map[string]interface{}

	// [TBD]: execute function

	return result, nil
}

func (integration JaegerIntegration) Validate() error {
	if integration.Config.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if integration.Config.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	return nil
}

func (integration JaegerIntegration) ValidateStep(
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
