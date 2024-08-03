package helpers

import (
	"encoding/json"
	"fmt"
)

func ValidateInputParameters(input any, parsedInput any) error {
	data, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("invalid input for get_properties_values function")
	}

	err = json.Unmarshal(data, &parsedInput)
	if err != nil {
		return fmt.Errorf("invalid input for get_properties_values function")
	}

	return nil
}
