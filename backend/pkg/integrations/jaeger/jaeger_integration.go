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
	"signal0ne/internal/utils"
	"signal0ne/pkg/integrations/helpers"
	"strings"
	"time"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"compare_traces": models.WorkflowFunctionDefinition{
		Function:   compareTraces,
		Input:      CompareTracesInput{},
		OutputTags: []string{"metadata"},
	},
	"get_properties_values": models.WorkflowFunctionDefinition{
		Function:   getPropertiesValues,
		Input:      GetPropertiesValuesInput{},
		OutputTags: []string{"logs", "traces"},
	},
	"get_dependencies": models.WorkflowFunctionDefinition{
		Function:   getDependencies,
		Input:      GetDependenciesInput{},
		OutputTags: []string{"metadata"},
	},
}

type JaegerIntegrationInventory struct {
	PyInterface net.Conn `json:"-" bson:"-"`
}

func NewJaegerIntegrationInventory(pyInterface net.Conn) JaegerIntegrationInventory {
	return JaegerIntegrationInventory{
		PyInterface: pyInterface,
	}
}

type JaegerIntegration struct {
	Inventory          JaegerIntegrationInventory
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:"config" bson:"config"`
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

	results = tools.ExecutionResultWrapper(intermediateResults, output, function.OutputTags)

	return results, nil
}

func (integration JaegerIntegration) Initialize() map[string]string {
	// Implement your config initialization here
	return nil
}

func (integration JaegerIntegration) Validate() error {
	if integration.Config.Url == "" {
		return fmt.Errorf("url cannot be empty")
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
	Service   string `json:"service"`
	Operation string `json:"operation"`

	TraceTags  string `json:"traceTags"`
	TraceQuery string `json:"traceQuery"`
}

type GetDependenciesInput struct {
	Service string `json:"service"`
}

func getPropertiesValues(input any, integration any) ([]any, error) {
	var parsedInput GetPropertiesValuesInput
	var output []any

	var finalComparisonField = make([]string, 0)

	err := helpers.ValidateInputParameters(input, &parsedInput, "get_properties_values")
	if err != nil {
		return output, err
	}

	comparedFieldParamSpliced := strings.Split(parsedInput.CompareBy, ",")
	for idx, field := range comparedFieldParamSpliced {
		comparedFieldParamSpliced[idx] = strings.Trim(field, " ")
	}

	fmt.Printf("###\nExecuting Jaeger integration function...\n")

	assertedIntegration := integration.(JaegerIntegration)

	url := assertedIntegration.Url
	apiPath := fmt.Sprintf("/api/traces?service=%s%s&tags=%s", parsedInput.Service, parsedInput.Query, parsedInput.Tags)

	finalUrl := fmt.Sprintf("%s%s", url, apiPath)

	intermediateTracesOutput, err := getJaegerObjects(finalUrl)
	if err != nil {
		return output, err
	}

	var serviceProcess string
	for _, trace := range intermediateTracesOutput {
		for key, process := range trace.(map[string]any)["processes"].(map[string]any) {
			if process.(map[string]any)["serviceName"] == parsedInput.Service {
				serviceProcess = key
				break
			}
		}
	}

	var spans []any
	for _, trace := range intermediateTracesOutput {
		for _, span := range trace.(map[string]any)["spans"].([]any) {
			assertedSpan, exists := span.(map[string]any)
			if !exists {
				continue
			}
			var spanWithDesiredValue = make(map[string]any)
			if assertedSpan["processID"] == serviceProcess {
				var intermediateSpan = make(map[string]any)
				intermediateSpan["logs"] = make([]map[string]any, 0)
				intermediateSpan["tags"] = make([]map[string]any, 0)

				for _, tag := range assertedSpan["tags"].([]any) {
					var parsedTag = make(map[string]any)
					assertedTag, exists := tag.(map[string]any)
					if !exists {
						continue
					}
					parsedTag[assertedTag["key"].(string)] = assertedTag["value"]
					intermediateSpan["tags"] = append(intermediateSpan["tags"].([]map[string]any), parsedTag)
				}

				for _, log := range assertedSpan["logs"].([]any) {
					var parsedLog = make(map[string]any)
					for _, field := range log.(map[string]any)["fields"].([]any) {
						assertedField, exists := field.(map[string]any)
						if !exists {
							continue
						}
						parsedLog[assertedField["key"].(string)] = assertedField["value"]
					}
					intermediateSpan["logs"] = append(intermediateSpan["logs"].([]map[string]any), parsedLog)
				}
				for _, comparisonField := range comparedFieldParamSpliced {
					fieldValuePlaceholder := tools.TraverseOutput(intermediateSpan, comparisonField, comparisonField)
					if fieldValuePlaceholder != nil || fieldValuePlaceholder != "" {
						spanWithDesiredValue[comparisonField] = fieldValuePlaceholder
						if !utils.Contains(finalComparisonField, comparisonField) {
							finalComparisonField = append(finalComparisonField, comparisonField)
						}
					}
				}
				if len(spanWithDesiredValue) == 0 {
					continue
				}
				spans = append(spans, spanWithDesiredValue)
			}
		}
	}

	pyInterfacePayload := map[string]any{
		"command": "get_log_occurrences",
		"params": map[string]any{
			"collectedLogs":  spans,
			"comparedFields": finalComparisonField,
		},
	}

	payloadBytes, err := json.Marshal(pyInterfacePayload)
	if err != nil {
		return output, err
	}

	batchSizeHeader := make([]byte, 4)
	binary.BigEndian.PutUint32(batchSizeHeader, uint32(len(payloadBytes)))
	payloadBytesWithHeaders := append(batchSizeHeader, payloadBytes...)

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

	for _, outputElement := range output {
		outputElement.(map[string]any)["output_source"] = parsedInput.Service
	}

	return output, nil

}

func compareTraces(input any, integration any) ([]any, error) {
	type Diff struct {
		Operation     string `json:"operation"`
		Spans         string `json:"spans"`
		DependencyMap string `json:"dependency_map"`
	}
	var parsedInput CompareTracesInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "compare_traces")
	if err != nil {
		return output, err
	}

	fmt.Printf("###\nExecuting Jaeger integration function...\n")

	assertedIntegration := integration.(JaegerIntegration)

	url := assertedIntegration.Url

	var operations []any
	var apiPath string
	var finalUrl string

	if parsedInput.Operation == "all" {
		apiPath = fmt.Sprintf("/api/services/%s/operations", parsedInput.Service)
		finalUrl = fmt.Sprintf("%s%s", url, apiPath)

		operations, err = getJaegerObjects(finalUrl)
		if err != nil {
			return output, err
		}
	} else {
		operations = []any{parsedInput.Operation}
	}

	for _, operation := range operations {
		var diff = Diff{}
		var traces []any

		tracesTags := strings.Split(parsedInput.TraceTags, ",")

		for _, tag := range tracesTags {
			apiPath = fmt.Sprintf("/api/traces?service=%s%s&operation=%s&limit=1&tags=%s", parsedInput.Service, parsedInput.TraceQuery, operation, tag)

			finalUrl = fmt.Sprintf("%s%s", url, apiPath)

			traces, err = getJaegerObjects(finalUrl)
			if err != nil {
				break
			}
			if len(traces) > 0 {
				break
			}
		}

		if len(traces) == 0 {
			fmt.Printf("No traces found for operation %s\n", operation)
			continue
		}

		baseProcesses := traces[0].(map[string]any)["processes"].(map[string]any)

		pid := 0
		diff.DependencyMap = fmt.Sprintf("%s\n", parsedInput.Service)
		for _, process := range baseProcesses {
			spacing := strings.Repeat("-", (pid+1)*2)
			diff.DependencyMap += fmt.Sprintf("%s%s\n", spacing, process.(map[string]any)["serviceName"].(string))
			pid++
		}

		diff.Operation = operation.(string)
		diff.Spans = "" //TODO: Implement spans comparison
		translatedMap := map[string]any{
			"dependency_map": diff.DependencyMap,
			"operation":      diff.Operation,
			"output_source":  parsedInput.Service,
		}

		output = append(output, translatedMap)

	}

	return output, nil

}

func getDependencies(input any, integration any) ([]any, error) {
	var parsedInput GetDependenciesInput
	var output []any
	var dependencyMap = make(map[string]any)
	var SERVICE_MAP_LOOKBACK = 604800000

	err := helpers.ValidateInputParameters(input, &parsedInput, "get_dependencies")
	if err != nil {
		return output, err
	}

	fmt.Printf("###\nExecuting Jaeger integration function...\n")

	assertedIntegration := integration.(JaegerIntegration)

	url := assertedIntegration.Url
	currentTimeInMs := time.Now().Unix() * 1000
	apiPath := fmt.Sprintf("/api/dependencies?endTs=%d&lookback=%d", currentTimeInMs, SERVICE_MAP_LOOKBACK)

	finalUrl := fmt.Sprintf("%s%s", url, apiPath)

	intermediateServiceOuput, err := getJaegerObjects(finalUrl)
	if err != nil {
		return output, fmt.Errorf("cannot retrieve service map, %s", err)
	}

	buildDependencyMap(intermediateServiceOuput, parsedInput.Service, dependencyMap)

	jsonifiedDependencyMap, err := json.Marshal(dependencyMap)
	if err != nil {
		return output, fmt.Errorf("cannot parse dependency map")
	}

	output = append(output, map[string]any{
		"dependency_map": string(jsonifiedDependencyMap),
	})

	return output, nil
}

func getJaegerObjects(url string) ([]any, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var bodyHandler map[string]any

	if resp.StatusCode != 200 {
		err = fmt.Errorf("%s", resp.Status)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &bodyHandler)
	if err != nil {
		err = fmt.Errorf("cannot parse response body, error %v", err)
		return nil, err
	}
	intermediateTracesOutput, exists := bodyHandler["data"].([]any)
	if !exists {
		err = fmt.Errorf("cannot parse response body")
		return nil, err
	}
	return intermediateTracesOutput, nil
}

func buildDependencyMap(servicesData []any, service string, serviceMap map[string]any) {
	for _, serviceObj := range servicesData {
		serviceName, ok := serviceObj.(map[string]any)["parent"].(string)
		if !ok {
			continue
		}
		if serviceName == service {
			child := serviceObj.(map[string]any)["child"].(string)
			serviceMap[service] = map[string]any{
				child: map[string]any{},
			}
			buildDependencyMap(servicesData, child, serviceMap[service].(map[string]any))
		}
	}
}
