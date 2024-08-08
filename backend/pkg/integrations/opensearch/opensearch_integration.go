package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
	"strconv"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/opensearch-project/opensearch-go/opensearchutil"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"get_log_occurrences": models.WorkflowFunctionDefinition{
		Function: getLogOccurrences,
		Input:    GetLogOccurrencesInput{},
	},
}

type OpenSearchIntegrationInventory struct {
	PyInterface net.Conn
}

func NewOpenSearchIntegrationInventory(pyInterface net.Conn) OpenSearchIntegrationInventory {
	return OpenSearchIntegrationInventory{
		PyInterface: pyInterface,
	}
}

type OpenSearchIntegration struct {
	Inventory          OpenSearchIntegrationInventory
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

func (integration OpenSearchIntegration) Execute(
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

func (integration OpenSearchIntegration) Validate() error {
	if integration.Config.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if integration.Config.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	return nil
}

func (integration OpenSearchIntegration) ValidateStep(
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

type GetLogOccurrencesInput struct {
	Query     string `json:"query"`
	CompareBy string `json:"compare_by"`
}

func getLogOccurrences(input any, integration any) ([]any, error) {
	var parsedInput GetLogOccurrencesInput
	var output []any
	var allLogObjects []any

	assertedIntegration := integration.(OpenSearchIntegration)

	err := helpers.ValidateInputParameters(input, &parsedInput, "get_log_occurrences")
	if err != nil {
		return output, err
	}

	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", assertedIntegration.Host, assertedIntegration.Port),
		},
	})
	if err != nil {
		return output, err
	}

	var query map[string]any
	err = json.Unmarshal([]byte(parsedInput.Query), &query)
	if err != nil {
		return output, err
	}

	searchReq := opensearchapi.SearchRequest{
		Index: []string{assertedIntegration.Index},
		Body:  opensearchutil.NewJSONReader(query),
	}

	searchResp, err := searchReq.Do(context.Background(), client)
	if err != nil {
		return output, fmt.Errorf("error performing search: %s", err)
	}
	defer searchResp.Body.Close()

	var hits map[string]any
	querySearchResults, err := io.ReadAll(searchResp.Body)
	if err != nil {
		return output, err
	}

	err = json.Unmarshal(querySearchResults, &hits)
	if err != nil {
		return output, err
	}
	parsedHits, ok := hits["hits"].(map[string]any)["hits"].([]any)
	if !ok {
		return output, fmt.Errorf("cannot parse output")
	}

	for _, hit := range parsedHits {
		intermediateHit, exists := hit.(map[string]any)["_source"]
		if !exists {
			return output, err
		}
		allLogObjects = append(allLogObjects, intermediateHit)
	}

	pyInterfacePayload := map[string]any{
		"command": "get_log_occurrences",
		"params": map[string]any{
			"collectedLogs":  allLogObjects,
			"comparedFields": parsedInput.CompareBy,
		},
	}
	payloadBytes, err := json.Marshal(pyInterfacePayload)
	if err != nil {
		return output, err
	}

	_, err = assertedIntegration.Inventory.PyInterface.Write(payloadBytes)
	if err != nil {
		return output, err
	}

	headerBuffer := make([]byte, 4)
	n, err := assertedIntegration.Inventory.PyInterface.Read(headerBuffer)
	if err != nil {
		return output, err
	}
	size, err := strconv.Atoi(string(headerBuffer[:n]))
	if err != nil {
		return output, err
	}

	payloadBuffer := make([]byte, size)
	n, err = assertedIntegration.Inventory.PyInterface.Read(payloadBuffer)
	if err != nil {
		return output, err
	}

	err = json.Unmarshal(headerBuffer[:n], &output)
	if err != nil {
		return output, err
	}

	return output, nil

}
