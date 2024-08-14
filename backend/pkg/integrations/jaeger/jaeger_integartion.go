package jaeger

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
	"strings"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"get_properties_values": models.WorkflowFunctionDefinition{
		Function: getPropertiesValues,
		Input:    GetPropertiesValuesInput{},
	},
	"compare_traces": models.WorkflowFunctionDefinition{
		Function: compareTraces,
		Input:    CompareTracesInput{},
	},
}

type JaegerIntegrationInventory struct {
	PyInterface net.Conn
}

func NewJaegerIntegrationInventory(pyInterface net.Conn) JaegerIntegrationInventory {
	return JaegerIntegrationInventory{
		PyInterface: pyInterface,
	}
}

type JaegerIntegration struct {
	Inventory          JaegerIntegrationInventory
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

func (integration JaegerIntegration) Execute(
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

type GetPropertiesValuesInput struct {
	Service   string `json:"service"`
	Tags      string `json:"tags"`
	Query     string `json:"query"`
	CompareBy string `json:"compare_by"`
}

type CompareTracesInput struct {
	Service        string `json:"service"`
	Operation      string `json:"operation"`
	Query          string `json:"query"`
	Tags           string `json:"tags"`
	ComparisonTags string `json:"comparisonTags"`
}

func getPropertiesValues(input any, integration any) ([]any, error) {
	var parsedInput GetPropertiesValuesInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "get_properties_values")
	if err != nil {
		return output, err
	}

	comparedFieldParamSpliced := strings.Split(parsedInput.CompareBy, ",")
	for idx, field := range comparedFieldParamSpliced {
		comparedFieldParamSpliced[idx] = strings.Trim(field, " ")
	}

	fmt.Printf("Executing Jaeger integration function...")

	assertedIntegration := integration.(JaegerIntegration)

	host := assertedIntegration.Host
	port := assertedIntegration.Port
	apiPath := fmt.Sprintf("/api/traces?service=%s%s&tags=%s", parsedInput.Service, parsedInput.Query, parsedInput.Tags)

	url := fmt.Sprintf("http://%s:%s%s", host, port, apiPath)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return output, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return output, err
	}
	defer resp.Body.Close()
	var bodyHandler map[string]any

	if resp.StatusCode != 200 {
		err = fmt.Errorf("%s", resp.Status)
		return output, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return output, err
	}

	err = json.Unmarshal(body, &bodyHandler)
	if err != nil {
		err = fmt.Errorf("cannot parse %s response body, error %v", assertedIntegration.Name, err)
		return output, err
	}
	spans, exists := bodyHandler["data"].([]any)
	if !exists {
		err = fmt.Errorf("cannot parse %s response body", assertedIntegration.Name)
		return output, err
	}
	fmt.Printf("No of spans:%d", len(spans))

	pyInterfacePayload := map[string]any{
		"command": "get_log_occurrences",
		"params": map[string]any{
			"collectedLogs":  spans,
			"comparedFields": comparedFieldParamSpliced,
		},
	}
	payloadBytes, err := json.Marshal(pyInterfacePayload)
	if err != nil {
		return output, err
	}

	headers := make([]byte, 4)
	binary.BigEndian.PutUint32(headers, uint32(len(payloadBytes)))
	payloadBytesWithHeaders := append(headers, payloadBytes...)
	fmt.Printf("SIZE: %s", string(headers))

	_, err = assertedIntegration.Inventory.PyInterface.Write(payloadBytesWithHeaders)
	if err != nil {
		return output, err
	}

	headerBuffer := make([]byte, 4)
	_, err = assertedIntegration.Inventory.PyInterface.Read(headerBuffer)
	if err != nil {
		return output, err
	}
	size := binary.BigEndian.Uint32(headerBuffer)

	payloadBuffer := make([]byte, size)
	n, err := assertedIntegration.Inventory.PyInterface.Read(payloadBuffer)
	if err != nil {
		return output, err
	}

	var intermediateOutput map[string]any
	err = json.Unmarshal(payloadBuffer[:n], &intermediateOutput)
	if err != nil {
		return output, err
	}
	statusCode, exists := intermediateOutput["status"].(string)
	if !exists || statusCode != "0" {
		errorMsg, _ := intermediateOutput["error"].(string)
		return output, fmt.Errorf("cannot retrieve results %s", errorMsg)
	}
	resultsEncoded, exists := intermediateOutput["result"].(string)
	if !exists {
		return output, fmt.Errorf("cannot retrieve results")
	}

	err = json.Unmarshal([]byte(resultsEncoded), &output)
	if err != nil {
		return output, err
	}

	return output, nil

}

func compareTraces(input any, integration any) ([]any, error) {
	var parsedInput CompareTracesInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "compare_traces")
	if err != nil {
		return output, err
	}

	fmt.Printf("Executing Jaeger integration function...")

	return output, nil

}
