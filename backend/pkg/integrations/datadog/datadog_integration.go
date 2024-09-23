package datadog

import (
	"context"
	"fmt"
	"signal0ne/internal/db"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"get_relevant_logs":    models.WorkflowFunctionDefinition{},
	"get_relevant_metrics": models.WorkflowFunctionDefinition{},
}

type DataDogIntegrationInventory struct {
	AlertsCollection *mongo.Collection `json:"-" bson:"-"`
}

func NewDataDogIntegrationInventory(
	alertsCollection *mongo.Collection) DataDogIntegrationInventory {
	return DataDogIntegrationInventory{
		AlertsCollection: alertsCollection,
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
