package confluence

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"search": models.WorkflowFunctionDefinition{
		Function: search,
		Input:    SearchInput{},
	},
}

type ConfluenceIntegrationInventory struct {
	PyInterface net.Conn `json:"-" bson:"-"`
}

func NewConfluenceIntegrationInventory(pyInterface net.Conn) ConfluenceIntegrationInventory {
	return ConfluenceIntegrationInventory{
		PyInterface: pyInterface,
	}
}

type ConfluenceIntegration struct {
	Inventory          ConfluenceIntegrationInventory
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
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
	Query          string `json:"query" bson:"query"`
	SimilarityCase string `json:"similarity_case" bson:"similarity_case"`
}

func search(input any, integration any) ([]any, error) {
	var parsedInput SearchInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "search")
	if err != nil {
		return output, err
	}

	assertedIntegration, _ := integration.(ConfluenceIntegration)

	//Hardcoded limit for UX readability reasons
	url := fmt.Sprintf("%s/wiki/rest/api/content/search?limit=6&expand=body.view&cql=%s", assertedIntegration.Url, parsedInput.Query)
	basicCredentials := base64.RawStdEncoding.EncodeToString([]byte(assertedIntegration.Email + ":" + assertedIntegration.ApiKey))

	contents, err := getPageContent(url, basicCredentials)
	if err != nil {
		return output, err
	}

	if parsedInput.SimilarityCase != "" {
		contents, err = assertedIntegration.compareContent(contents, parsedInput.SimilarityCase)
		if err != nil {
			return output, err
		}
	}

	for _, content := range contents {
		output = append(output, content)
	}

	return output, nil
}

func (integration ConfluenceIntegration) compareContent(contents []string, similarityCase string) ([]string, error) {
	var similarContents []string

	pyInterfacePayload := map[string]any{
		"command": "contents_similarity",
		"params": map[string]any{
			"similarityCase": similarityCase,
			"contents":       contents,
		},
	}

	payloadBytes, err := json.Marshal(pyInterfacePayload)
	if err != nil {
		return similarContents, err
	}

	batchSizeHeader := make([]byte, 4)
	binary.BigEndian.PutUint32(batchSizeHeader, uint32(len(payloadBytes)))
	payloadBytesWithHeaders := append(batchSizeHeader, payloadBytes...)

	_, err = integration.Inventory.PyInterface.Write(payloadBytesWithHeaders)
	if err != nil {
		return similarContents, err
	}

	headerBuffer := make([]byte, 4)
	_, err = integration.Inventory.PyInterface.Read(headerBuffer)
	if err != nil {
		return similarContents, err
	}

	size := binary.BigEndian.Uint32(headerBuffer)

	payloadBuffer := make([]byte, size)
	numberOfBytesSent, err := integration.Inventory.PyInterface.Read(payloadBuffer)
	if err != nil {
		return similarContents, err
	}

	var intermediateOutput map[string]any
	err = json.Unmarshal(payloadBuffer[:numberOfBytesSent], &intermediateOutput)
	if err != nil {
		return similarContents, err
	}

	statusCode, exists := intermediateOutput["status"].(string)
	if !exists || statusCode != "0" {
		errorMsg, _ := intermediateOutput["error"].(string)
		return similarContents, fmt.Errorf("cannot retrieve results %s", errorMsg)
	}

	resultsEncoded, exists := intermediateOutput["result"].(string)
	if !exists {
		return similarContents, fmt.Errorf("cannot retrieve results")
	}

	err = json.Unmarshal([]byte(resultsEncoded), &similarContents)
	if err != nil {
		return similarContents, err
	}

	return similarContents, nil
}

func getPageContent(url string, credentials string) ([]string, error) {
	var contents []string
	var results []any

	fmt.Printf("Getting content from confluence %s\n", url)

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
		return contents, fmt.Errorf("failed to get content from confluence %s", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return contents, err
	}

	var interfaceResults map[string]any
	err = json.Unmarshal(responseBody, &interfaceResults)
	if err != nil {
		return contents, err
	}

	results = interfaceResults["results"].([]any)

	for _, result := range results {
		content := result.(map[string]any)["body"].(map[string]any)["view"].(map[string]any)["value"].(string)
		contents = append(contents, content)
	}

	return contents, nil
}
