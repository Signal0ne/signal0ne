package alertmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	"get_relevant_alerts": models.WorkflowFunctionDefinition{
		Function: getRelevantAlerts,
		Input:    GetRelevantAlertsInput{},
	},
}

type AlertmanagerIntegrationInventory struct {
	AlertsCollection *mongo.Collection `json:"-" bson:"-"`
}

func NewAlertmanagerIntegrationInventory(
	alertsCollection *mongo.Collection) AlertmanagerIntegrationInventory {
	return AlertmanagerIntegrationInventory{
		AlertsCollection: alertsCollection,
	}
}

type AlertmanagerIntegration struct {
	Config             `json:"config" bson:"config"`
	Inventory          AlertmanagerIntegrationInventory
	models.Integration `json:",inline" bson:",inline"`
}

func (integration AlertmanagerIntegration) Trigger(
	payload map[string]any,
	alert *models.EnrichedAlert,
	workflow *models.Workflow) (err error) {

	var StateKey = "status"

	//incoming in RFC3339 format string with timezone UTC
	var StartTimeKey = "startsAt"

	//TBD: WE DO NOT SUPPORT ALERT GROUPING FOR ALERTMANAGER
	//supported alertmanager config signal0ne receiver:
	//```
	//route:
	//   receiver: "singal0ne"
	//   group_by: ['...']
	//```
	alertPayload, exists := payload["alerts"].([]any)[0].(map[string]any)
	if !exists {
		return fmt.Errorf("cannot find alerts in payload")
	}

	alert.TriggerProperties, err = tools.WebhookTriggerExec(alertPayload, workflow)
	if err != nil {
		return err
	}

	alert.State, err = tools.MapAlertState(alertPayload, StateKey, TriggerStateMapping)
	if err != nil {
		return err
	}

	alert.StartTime, err = tools.GetStartTime(alertPayload, StartTimeKey)
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

func (integration AlertmanagerIntegration) Execute(
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

func (integration AlertmanagerIntegration) Validate() error {
	if integration.Config.Url == "" {
		return fmt.Errorf("url cannot be empty")
	}

	return nil
}

func (integration AlertmanagerIntegration) ValidateStep(
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

type GetRelevantAlertsInput struct {
	Filter string `json:"filter"`
}

func getRelevantAlerts(input any, integration any) ([]any, error) {
	var parsedInput GetRelevantAlertsInput
	var output []any

	err := helpers.ValidateInputParameters(input, &parsedInput, "compare_traces")
	if err != nil {
		return output, err
	}

	fmt.Printf("###\nExecuting Alertmanager integration function...\n")

	assertedIntegration := integration.(AlertmanagerIntegration)
	apiPath := "/api/v2/alerts?"

	url := fmt.Sprintf("%s%s", assertedIntegration.Url, apiPath)

	filters := strings.Split(parsedInput.Filter, ",")

	alerts, err := getAlerts(url, filters)
	if err != nil {
		return output, err
	}

	for _, alert := range alerts {
		source, exists := alert.(map[string]any)["labels"].(map[string]any)["name"].(string)
		if !exists {
			source = ""
		}
		output = append(output, map[string]any{
			"alert":         alert,
			"output_source": source,
		})
	}

	return output, nil
}

func getAlerts(url string, filters []string) ([]any, error) {
	client := &http.Client{}

	if !(len(filters) > 0) {
		return nil, nil
	} else {
		for fi, filter := range filters {
			if filter != "" {
				if fi != 0 {
					url += "&"
				}
				url += ("filter=" + filter)
			}
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var bodyHandler []any

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

	return bodyHandler, nil
}
