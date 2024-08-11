package signal0ne

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"correlate_ongoing_alerts": models.WorkflowFunctionDefinition{
		Function: correlateOngoingAlerts,
		Input:    CorrelateOngoingAlertsInput{},
	},
}

type Signal0neIntegrationInventory struct {
	AlertsCollection *mongo.Collection
	PyInterface      net.Conn
}

func NewOpenSearchIntegrationInventory(
	alertsCollection *mongo.Collection,
	pyInterface net.Conn) Signal0neIntegrationInventory {
	return Signal0neIntegrationInventory{
		AlertsCollection: alertsCollection,
		PyInterface:      pyInterface,
	}
}

type Signal0neIntegration struct {
	Inventory          Signal0neIntegrationInventory
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:",inline" bson:",inline"`
}

func (integration Signal0neIntegration) Execute(
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

func (integration Signal0neIntegration) Validate() error {
	// Implement your config validation here
	return nil
}

func (integration Signal0neIntegration) ValidateStep(
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

type CorrelateOngoingAlertsInput struct {
	StartTimestamp string `json:"startTimestamp"`
	EndTimestamp   string `json:"endTimestamp"`
	Type           string `json:"type"`
	CompareBy      string `json:"compareBy"`
}

func correlateOngoingAlerts(input any, integration any) ([]any, error) {
	var parsedInput CorrelateOngoingAlertsInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "correlate_ongoing_alerts")
	if err != nil {
		return output, err
	}

	assertedIntegration := integration.(Signal0neIntegration)

	comparedFieldParamSpliced := strings.Split(parsedInput.CompareBy, ",")
	for idx, field := range comparedFieldParamSpliced {
		comparedFieldParamSpliced[idx] = strings.Trim(field, " ")
	}

	fmt.Printf("Executing backstage getPropertiesValues\n")

	unixStartTimestamp, err := strconv.Atoi(parsedInput.StartTimestamp)
	if err != nil {
		return output, err
	}
	unixEndTimestamp, err := strconv.Atoi(parsedInput.EndTimestamp)
	if err != nil {
		return output, err
	}

	startTimestamp := time.Unix(int64(unixStartTimestamp), 0)
	endTimestamp := time.Unix(int64(unixEndTimestamp), 0)

	filter := bson.M{
		"timestamp": bson.M{
			"$gte": startTimestamp,
			"$lte": endTimestamp,
		},
	}

	var alerts []models.EnrichedAlert
	potentialCorrelationsResults, err := assertedIntegration.Inventory.AlertsCollection.Find(context.Background(), filter)
	if err != nil {
		return output, err
	}
	potentialCorrelationsResults.Decode(&alerts)

	pyInterfacePayload := map[string]any{
		"command": "correlate_ongoing_alerts",
		"params": map[string]any{
			"collectedAlerts": alerts,
			"comparedFields":  comparedFieldParamSpliced,
		},
	}
	payloadBytes, err := json.Marshal(pyInterfacePayload)
	if err != nil {
		return output, err
	}

	headers := make([]byte, 4)
	binary.BigEndian.PutUint32(headers, uint32(len(payloadBytes)))
	payloadBytesWithHeaders := append(headers, payloadBytes...)

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
