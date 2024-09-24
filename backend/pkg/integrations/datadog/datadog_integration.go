package datadog

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"signal0ne/internal/db"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"get_relevant_logs": models.WorkflowFunctionDefinition{
		Function: getRelevantLogs,
		Input:    GetRelevantLogsInput{},
	},

	// TBD
	"get_triggered_monitors": models.WorkflowFunctionDefinition{},
	"create_incident":        models.WorkflowFunctionDefinition{},
}

type DataDogIntegrationInventory struct {
	AlertsCollection   *mongo.Collection `json:"-" bson:"-"`
	WorkflowProperties *models.Workflow  `json:"-" bson:"-"`
	PyInterface        net.Conn          `json:"-" bson:"-"`
}

func NewDataDogIntegrationInventory(
	pyInterface net.Conn,
	alertsCollection *mongo.Collection,
	workflowProperties *models.Workflow) DataDogIntegrationInventory {
	return DataDogIntegrationInventory{
		AlertsCollection:   alertsCollection,
		WorkflowProperties: workflowProperties,
		PyInterface:        pyInterface,
	}
}

type DataDogIntegration struct {
	Config             `json:"config" bson:"config"`
	Inventory          DataDogIntegrationInventory
	models.Integration `json:",inline" bson:",inline"`
}

func (integration DataDogIntegration) Trigger(
	payload map[string]any,
	alert *models.EnrichedAlert,
	workflow *models.Workflow) (err error) {

	var StateKey = "ALERT_TRANSITION"

	//incoming in unix timestamp format
	var StartTimeKey = "DATE"

	alert.TriggerProperties, err = tools.WebhookTriggerExec(payload, workflow)
	if err != nil {
		return err
	}

	alert.State, err = tools.MapAlertState(payload, StateKey, TriggerStateMapping)
	if err != nil {
		return err
	}

	alert.StartTime, err = tools.GetStartTime(payload, StartTimeKey)
	if err != nil {
		return err
	}

	alertsHistory, err := db.GetEnrichedAlertsByWorkflowId(workflow.Id.Hex(),
		context.Background(),
		integration.Inventory.AlertsCollection,
		bson.M{
			"startTime": alert.StartTime,
		},
	)
	if err != nil {
		return err
	}

	if len(alertsHistory) > 0 {
		var anyUpdates = false
		for _, alertFromHistory := range alertsHistory {
			alertFromHistory.State = alert.State
			err = db.UpdateEnrichedAlert(alertFromHistory, context.Background(), integration.Inventory.AlertsCollection)
			if err != nil {
				continue
			} else {
				anyUpdates = true
			}
		}

		if anyUpdates && alert.State == models.AlertStatusInactive {
			return fmt.Errorf("alert already inactive")
		}
	}

	return nil
}

func (integration DataDogIntegration) Execute(
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

func (integration DataDogIntegration) Validate() error {
	if integration.Config.Url == "" {
		return fmt.Errorf("url cannot be empty")
	}

	return nil
}

func (integration DataDogIntegration) ValidateStep(
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

type GetRelevantLogsInput struct {
	CompareBy string `json:"compare_by" bson:"compare_by"`
	Index     string `json:"index" bson:"index"`
	Limit     int    `json:"limit" bson:"limit"`
	Query     string `json:"query" bson:"query"`
	Service   string `json:"service" bson:"service"`
}

func getRelevantLogs(input any, integration any) ([]any, error) {
	var parsedInput GetRelevantLogsInput
	var output []any
	var query string
	var allLogObjects []any

	assertedIntegration := integration.(DataDogIntegration)

	err := helpers.ValidateInputParameters(input, &parsedInput, "get_relevant_logs")
	if err != nil {
		return output, err
	}

	fmt.Printf("###\nExecuting OpenSearch integration function...\n")

	comparedFieldParamSpliced := strings.Split(parsedInput.CompareBy, ",")
	for idx, field := range comparedFieldParamSpliced {
		comparedFieldParamSpliced[idx] = strings.Trim(field, " ")
	}

	client := &http.Client{}

	url := fmt.Sprintf("%s/api/v2/logs/events/search", assertedIntegration.Config.Url)

	if parsedInput.Index == "" {
		parsedInput.Index = "*"
	}

	indexes := []string{parsedInput.Index}
	from := fmt.Sprintf("now-%s", assertedIntegration.Inventory.WorkflowProperties.Lookback)
	to := "now"

	if parsedInput.Service != "" {
		query = fmt.Sprintf("service:%s %s", parsedInput.Service, parsedInput.Query)
	} else {
		query = parsedInput.Query
	}

	requestBody := map[string]any{
		"filter": map[string]any{
			"query":   query,
			"from":    from,
			"to":      to,
			"indexes": indexes,
		},
		"page": map[string]any{
			"limit": parsedInput.Limit,
		},
		"sort": "timestamp",
	}

	encodedRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(encodedRequestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", assertedIntegration.ApiKey)
	req.Header.Set("DD-APPLICATION-KEY", assertedIntegration.ApplicationKey)
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

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &bodyHandler)
	if err != nil {
		err = fmt.Errorf("cannot parse response body, error %v", err)
		return nil, err
	}
	intermediateLogsOutput, exists := bodyHandler["data"].([]any)
	if !exists {
		err = fmt.Errorf("cannot parse response body")
		return nil, err
	}

	if len(intermediateLogsOutput) == 0 {
		return output, nil
	}

	for _, hit := range intermediateLogsOutput {
		intermediateHit, exists := hit.(map[string]any)
		var parsedIntermediateHit = make(map[string]any)
		if !exists {
			return output, err
		}
		for _, mapping := range comparedFieldParamSpliced {
			parsedIntermediateHit[mapping] = tools.TraverseOutput(intermediateHit, mapping, mapping)
		}
		allLogObjects = append(allLogObjects, parsedIntermediateHit)
	}

	pyInterfacePayload := map[string]any{
		"command": "get_log_occurrences",
		"params": map[string]any{
			"collectedLogs":  allLogObjects,
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
		return output, fmt.Errorf("cannot retrieve results, results not found")
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
