package signal0ne

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"

	"go.mongodb.org/mongo-driver/mongo"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"correlate_ongoing_alerts": models.WorkflowFunctionDefinition{
		Function: correlateOngoingAlerts,
		Input:    CorrelateOngoingAlertsInput{},
	},
	"create_incident": models.WorkflowFunctionDefinition{
		Function: createIncident,
		Input:    CreateIncidentInput{},
	},
}

type Signal0neIntegrationInventory struct {
	IncidentCollection *mongo.Collection `json:"-" bson:"-"`
	PyInterface        net.Conn          `json:"-" bson:"-"`
	WorkflowProperties *models.Workflow  `json:"-" bson:"-"`
}

func NewSignal0neIntegrationInventory(
	incidentsCollection *mongo.Collection,
	pyInterface net.Conn,
	workflowProperties *models.Workflow) Signal0neIntegrationInventory {
	return Signal0neIntegrationInventory{
		IncidentCollection: incidentsCollection,
		PyInterface:        pyInterface,
		WorkflowProperties: workflowProperties,
	}
}

type Signal0neIntegration struct {
	Inventory          Signal0neIntegrationInventory
	models.Integration `json:",inline" bson:",inline"`
	Config             `json:"config" bson:"config"`
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

type CreateIncidentInput struct {
	Severity              string `json:"severity"`
	Assignee              string `json:"assignee"`
	ParsableContextObject string `json:"parsable_context_object"`
}

func correlateOngoingAlerts(input any, integration any) ([]any, error) {
	var parsedInput CorrelateOngoingAlertsInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "correlate_ongoing_alerts")
	if err != nil {
		return output, err
	}

	// assertedIntegration := integration.(Signal0neIntegration)

	// comparedFieldParamSpliced := strings.Split(parsedInput.CompareBy, ",")
	// for idx, field := range comparedFieldParamSpliced {
	// 	comparedFieldParamSpliced[idx] = strings.Trim(field, " ")
	// }

	// fmt.Printf("Executing backstage getPropertiesValues\n")

	// unixStartTimestamp, err := strconv.Atoi(parsedInput.StartTimestamp)
	// if err != nil {
	// 	return output, err
	// }
	// unixEndTimestamp, err := strconv.Atoi(parsedInput.EndTimestamp)
	// if err != nil {
	// 	return output, err
	// }

	// startTimestamp := time.Unix(int64(unixStartTimestamp), 0)
	// endTimestamp := time.Unix(int64(unixEndTimestamp), 0)

	// filter := bson.M{
	// 	"timestamp": bson.M{
	// 		"$gte": startTimestamp,
	// 		"$lte": endTimestamp,
	// 	},
	// }

	// var alerts []models.EnrichedAlert
	// potentialCorrelationsResults, err := assertedIntegration.Inventory.AlertsCollection.Find(context.Background(), filter)
	// if err != nil {
	// 	return output, err
	// }
	// potentialCorrelationsResults.Decode(&alerts)

	// var entities = make([]any, 0)
	// for _, _ = range alerts {
	// 	// TBD
	// }

	// pyInterfacePayload := map[string]any{
	// 	"command": "correlate_ongoing_alerts",
	// 	"params": map[string]any{
	// 		"collectedEntities": entities,
	// 	},
	// }
	// payloadBytes, err := json.Marshal(pyInterfacePayload)
	// if err != nil {
	// 	return output, err
	// }

	// headers := make([]byte, 4)
	// binary.BigEndian.PutUint32(headers, uint32(len(payloadBytes)))
	// payloadBytesWithHeaders := append(headers, payloadBytes...)

	// _, err = assertedIntegration.Inventory.PyInterface.Write(payloadBytesWithHeaders)
	// if err != nil {
	// 	return output, err
	// }
	// headerBuffer := make([]byte, 4)
	// _, err = assertedIntegration.Inventory.PyInterface.Read(headerBuffer)
	// if err != nil {
	// 	return output, err
	// }
	// size := binary.BigEndian.Uint32(headerBuffer)

	// payloadBuffer := make([]byte, size)
	// n, err := assertedIntegration.Inventory.PyInterface.Read(payloadBuffer)
	// if err != nil {
	// 	return output, err
	// }

	// var intermediateOutput map[string]any
	// err = json.Unmarshal(payloadBuffer[:n], &intermediateOutput)
	// if err != nil {
	// 	return output, err
	// }
	// statusCode, exists := intermediateOutput["status"].(string)
	// if !exists || statusCode != "0" {
	// 	errorMsg, _ := intermediateOutput["error"].(string)
	// 	return output, fmt.Errorf("cannot retrieve results %s", errorMsg)
	// }
	// resultsEncoded, exists := intermediateOutput["result"].(string)
	// if !exists {
	// 	return output, fmt.Errorf("cannot retrieve results")
	// }

	// err = json.Unmarshal([]byte(resultsEncoded), &output)
	// if err != nil {
	// 	return output, err
	// }

	return output, nil
}

func createIncident(input any, integration any) ([]any, error) {
	var parsedInput CreateIncidentInput
	var parsedAlert models.EnrichedAlert
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "create_incident")
	if err != nil {
		return output, err
	}

	assertedIntegration := integration.(Signal0neIntegration)

	err = json.Unmarshal([]byte(parsedInput.ParsableContextObject), &parsedAlert)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Get primary fields
	primaryFields := make([]map[string]any, 0)
	primaryFields = append(primaryFields, parsedAlert.TriggerProperties)

	// Search for assignee in db
	assignee := models.User{}

	// Filling in Tasks with Items
	tasks := make([]models.Task, 0)

	for si, step := range assertedIntegration.Inventory.WorkflowProperties.Steps {
		isDone := true

		stepOutputs := parsedAlert.AdditionalContext[fmt.Sprint("%s_%s", step.Integration, step.Function)].Output.([]any)

		// Check if is done
		if len(stepOutputs) == 0 {
			isDone = false
		}

		task := models.Task{
			StepName: step.Name,
			Priority: si,
			Assignee: models.User{},
			IsDone:   isDone,
			Items:    make([]models.Item, 0),
		}

		for _, stepOutput := range stepOutputs {
			item := models.Item{
				Fields: make([]models.Field, 0),
			}
			for key, value := range stepOutput.(map[string]any) {
				item.Fields = append(item.Fields, models.Field{
					Key:       key,
					Source:    step.Integration,
					Value:     value,
					ValueType: "markdown",
				})
			}
			task.Items = append(task.Items, item)
		}
		tasks = append(tasks, task)
	}

	incident := models.Incident{
		Id:            parsedAlert.Id,
		Title:         assertedIntegration.Inventory.WorkflowProperties.Name,
		Assignee:      assignee,
		Severity:      parsedInput.Severity,
		PrimaryFields: primaryFields,
		Tasks:         tasks,
		History:       []models.IncidentUpdate[models.Update]{},
	}

	assertedIntegration.Inventory.IncidentCollection.InsertOne(context.Background(), incident)
	return output, nil
}
