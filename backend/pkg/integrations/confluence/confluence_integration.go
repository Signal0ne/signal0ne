package confluence

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
)

type ConfluenceIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

var functions = map[string]models.WorkflowFunctionDefinition{
	"search": models.WorkflowFunctionDefinition{
		Function: search,
		Input:    SearchInput{},
	},
}

func (integration ConfluenceIntegration) Execute(
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

func (integration ConfluenceIntegration) Validate() error {
	// Implement your config validation here
	return nil
}

func (integration ConfluenceIntegration) ValidateStep(
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

type SearchInput struct {
	Query string `json:"query" bson:"query"`
}

func search(input any, integration any) ([]any, error) {
	var parsedInput SearchInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "get_content")
	if err != nil {
		return output, err
	}

	assertedIntegration, _ := integration.(ConfluenceIntegration)

	//Hardcoded limit for UX readability reasons
	url := fmt.Sprintf("%s/rest/api/content?limit=6&cql=%s", assertedIntegration.Url, parsedInput.Query)
	basicCredentials := base64.RawStdEncoding.EncodeToString([]byte(assertedIntegration.Email + ":" + assertedIntegration.APIKey))

	contents, err := getPageContent(url, basicCredentials)
	if err != nil {
		return output, err
	}

	for _, content := range contents {
		output = append(output, content)
	}

	return output, nil

}

func getPageContent(url string, credentials string) ([]string, error) {
	var contents []string
	var results []any

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return contents, err
	}

	req.Header.Set("Authorization", "Basic "+credentials)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return contents, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return contents, fmt.Errorf("failed to get content from confluence")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return contents, err
	}

	err = json.Unmarshal(responseBody, &results)
	if err != nil {
		return contents, err
	}

	for _, result := range results {
		content := result.(map[string]any)["body"].(map[string]any)["view"].(map[string]any)["value"].(string)
		contents = append(contents, content)
	}

	return contents, nil
}
