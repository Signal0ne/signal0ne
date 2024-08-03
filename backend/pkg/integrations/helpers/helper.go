package helpers

import (
	"encoding/json"
	"fmt"
)

func ValidateInputParameters(input any, parsedInput any, functionName string) error {
	data, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("invalid input for %s function")
	}

	err = json.Unmarshal(data, &parsedInput)
	if err != nil {
		return fmt.Errorf("invalid input for %s function")
	}

	return nil
}
