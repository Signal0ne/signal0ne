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

func (i JaegerIntegration) Execute(
	input interface{},
	output interface{},
	functionName string) (map[string]interface{}, error) {

	var result map[string]interface{}

	// [TBD]: execute funtion

	return result, nil
}

func (i JaegerIntegration) Validate() error {
	if i.Config.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if i.Config.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	return nil
}

func (i JaegerIntegration) ValidateStep(
	input interface{},
	output interface{},
	functionName string,
) error {
	function, exists := functions[functionName]
	if !exists {
		return fmt.Errorf("cannot find selected funtion")
	}

	_, err := function(input, true)
	if err != nil {
		return err
	}

	return nil
}
