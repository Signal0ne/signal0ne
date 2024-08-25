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
	Service   string `json:"service"`
	Operation string `json:"operation"`

	BaseTraceTags      string `json:"baseTraceTags"`
	ComparedTraceTags  string `json:"comparedTraceTags"`
	BaseTraceQuery     string `json:"baseTraceQuery"`
	ComparedTraceQuery string `json:"comparedTraceQuery"`
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

	fmt.Printf("Executing Jaeger integration function...\n")

	assertedIntegration := integration.(JaegerIntegration)

	host := assertedIntegration.Host
	port := assertedIntegration.Port
	apiPath := fmt.Sprintf("/api/traces?service=%s%s&tags=%s", parsedInput.Service, parsedInput.Query, parsedInput.Tags)

	url := fmt.Sprintf("http://%s:%s%s", host, port, apiPath)

	intermediateTracesOutput, err := getJaegerObjects(url)
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
			assertedSpan := span.(map[string]any)
			var spanWithDesiredValue = make(map[string]any)
			if assertedSpan["processID"] == serviceProcess {
				var intermediateSpan = make(map[string]any)
				intermediateSpan["logs"] = make([]map[string]any, 0)
				intermediateSpan["tags"] = make([]map[string]any, 0)

				for _, tag := range assertedSpan["tags"].([]any) {
					var parsedTag = make(map[string]any)
					assertedTag := tag.(map[string]any)
					parsedTag[assertedTag["key"].(string)] = assertedTag["value"]
					intermediateSpan["tags"] = append(intermediateSpan["tags"].([]map[string]any), parsedTag)
				}

				for _, log := range assertedSpan["logs"].([]any) {
					var parsedLog = make(map[string]any)
					for _, field := range log.(map[string]any)["fields"].([]any) {
						assertedField := field.(map[string]any)
						parsedLog[assertedField["key"].(string)] = assertedField["value"]
					}
					intermediateSpan["logs"] = append(intermediateSpan["logs"].([]map[string]any), parsedLog)
				}
				for _, comparisonField := range comparedFieldParamSpliced {
					spanWithDesiredValue[comparisonField] = tools.TraverseOutput(intermediateSpan, comparisonField, comparisonField)
				}
				spans = append(spans, spanWithDesiredValue)
			}
		}
	}

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
		Operation string         `json:"operation"`
		Processes map[string]any `json:"processes"`
		Spans     map[string]any `json:"spans"`
	}
	var parsedInput CompareTracesInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "compare_traces")
	if err != nil {
		return output, err
	}

	fmt.Printf("Executing Jaeger integration function...")

	assertedIntegration := integration.(JaegerIntegration)

	host := assertedIntegration.Host
	port := assertedIntegration.Port

	var operations []any
	var apiPath string
	var url string
	if parsedInput.Operation == "all" {
		apiPath = fmt.Sprintf("/api/services/%s/operations", parsedInput.Service)
		url = fmt.Sprintf("http://%s:%s%s", host, port, apiPath)

		operations, err = getJaegerObjects(url)
		if err != nil {
			return output, err
		}
	} else {
		operations = []any{parsedInput.Operation}
	}

	for _, operation := range operations {
		var diff = Diff{}
		var traces []any
		var tracesToCompare []any

		tracesTags := strings.Split(parsedInput.BaseTraceTags, ",")
		tracesToCompareTags := strings.Split(parsedInput.ComparedTraceTags, ",")

		//BaseTraces
		for _, tag := range tracesTags {
			apiPath = fmt.Sprintf("/api/traces?service=%s%s&operation=%s&limit=1&tags=%s", parsedInput.Service, parsedInput.BaseTraceQuery, operation, tag)

			url = fmt.Sprintf("http://%s:%s%s", host, port, apiPath)

			traces, err = getJaegerObjects(url)
			if err != nil {
				break
			}
			if len(traces) > 0 {
				break
			}
		}

		//ComparedTraces
		for _, tag := range tracesToCompareTags {
			apiPath = fmt.Sprintf("/api/traces?service=%s%s&operation=%s&limit=1&tags=%s", parsedInput.Service, parsedInput.ComparedTraceQuery, operation, tag)

			url = fmt.Sprintf("http://%s:%s%s", host, port, apiPath)

			tracesToCompare, err = getJaegerObjects(url)
			if err != nil {
				break
			}
			if len(tracesToCompare) > 0 {
				break
			}
		}

		if len(traces) == 0 || len(tracesToCompare) == 0 {
			continue
		}

		//Compare processes
		baseProcesses := traces[0].(map[string]any)["processes"].(map[string]any)
		comparedProcesses := tracesToCompare[0].(map[string]any)["processes"].(map[string]any)

		baseProcessesSlice := make([]string, 0)
		comparedProcessesSlice := make([]string, 0)

		for _, process := range baseProcesses {
			baseProcessesSlice = append(baseProcessesSlice, process.(map[string]any)["serviceName"].(string))
		}

		for _, process := range comparedProcesses {
			comparedProcessesSlice = append(comparedProcessesSlice, process.(map[string]any)["serviceName"].(string))
		}

		processesDiffSlice := diffStringSlices(baseProcessesSlice, comparedProcessesSlice)
		//Compare spans, errors, durations --- TBD
		if len(processesDiffSlice) > 0 {
			for _, process := range processesDiffSlice {
				sign := string(process[0])
				processName := string(process[1:])
				diff.Processes = map[string]any{
					"sign":        sign,
					"processName": processName,
				}
			}

			diff.Operation = operation.(string)
			translatedMap := map[string]any{
				"output_source": parsedInput.Service,
				"processes":     diff.Processes,
				"operation":     diff.Operation,
			}
			output = append(output, translatedMap)
		}
	}

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

func diffStringSlices(old, new []string) []string {
	var diff = make([]string, 0)
	i, j, k := 0, 0, 0
	for ; i < len(old); i++ {
		present := false
		j = 0
		k = 0

		// Skip if already exists in diff
		for ; k < len(diff); k++ {
			if old[i] == string(diff[k][1:]) {
				present = true
				break
			}
		}
		if present {
			continue
		}

		for ; j < len(new); j++ {
			if old[i] == new[j] {
				present = true
				break
			}
		}
		if !present {
			diff = append(diff, "-"+old[i])
		}
	}

	return diff
}
