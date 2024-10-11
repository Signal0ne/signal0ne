package signal0ne

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"signal0ne/cmd/config"
	"signal0ne/internal/db"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations/helpers"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var functions = map[string]models.WorkflowFunctionDefinition{
	"correlate_ongoing_alerts": models.WorkflowFunctionDefinition{
		Function:   correlateOngoingAlerts,
		Input:      CorrelateOngoingAlertsInput{},
		OutputTags: []string{"alerts"},
	},
	"create_incident": models.WorkflowFunctionDefinition{
		Function:   createIncident,
		Input:      CreateIncidentInput{},
		OutputTags: []string{"metadata"},
	},
}

type Signal0neIntegrationInventory struct {
	AlertsCollection   *mongo.Collection `json:"-" bson:"-"`
	IncidentCollection *mongo.Collection `json:"-" bson:"-"`
	PyInterface        net.Conn          `json:"-" bson:"-"`
	WorkflowProperties *models.Workflow  `json:"-" bson:"-"`
}

func NewSignal0neIntegrationInventory(
	alertsCollection *mongo.Collection,
	incidentsCollection *mongo.Collection,
	pyInterface net.Conn,
	workflowProperties *models.Workflow) Signal0neIntegrationInventory {
	return Signal0neIntegrationInventory{
		AlertsCollection:   alertsCollection,
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

	results = tools.ExecutionResultWrapper(intermediateResults, output, function.OutputTags)

	return results, nil
}

func (integration Signal0neIntegration) Initialize() map[string]string {
	// Implement your config initialization here
	return nil
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
	DependencyMap  string `json:"dependency_map"`
}

type CreateIncidentInput struct {
	Severity                 string   `json:"severity"`
	Assignee                 string   `json:"assignee"`
	ParsableContextObject    string   `json:"parsable_context_object"`
	ManuallyCorrelatedAlerts []string `json:"_manually_correlated_alerts"`
}

func correlateOngoingAlerts(input any, integration any) ([]any, error) {
	var parsedInput CorrelateOngoingAlertsInput
	var output []any
	var services []string

	err := helpers.ValidateInputParameters(input, &parsedInput, "correlate_ongoing_alerts")
	if err != nil {
		return output, err
	}

	parsedIntegration := integration.(Signal0neIntegration)

	serviceDependencyMap := strings.Split(parsedInput.DependencyMap, ",")

	for _, serviceDependency := range serviceDependencyMap {
		operationServiceMap := strings.Split(serviceDependency, "\n")[1:]
		for si, service := range operationServiceMap {
			prefix := strings.Repeat("-", (si+1)*2)
			services = append(services, strings.TrimPrefix(service, prefix))
		}
	}

	startTime, err := time.Parse(time.RFC3339, parsedInput.StartTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start time: %v", err)
	}

	lookback := parsedIntegration.Inventory.WorkflowProperties.Lookback
	lookbackSuffix := lookback[len(lookback)-1:]
	lookbackValue, err := strconv.Atoi(lookback[:len(lookback)-1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse lookback: %v", err)
	}

	delta := -time.Duration(lookbackValue)

	if lookbackSuffix == "m" {
		delta = delta * time.Minute
	} else if lookbackSuffix == "s" {
		delta = delta * time.Second
	}

	endTime := startTime.Add(delta)

	q := bson.M{
		"startTime": bson.M{
			"$gte": endTime,
			"$lte": startTime,
		},
		"service": bson.M{
			"$in": services,
		},
		"state": models.AlertStatusActive,
	}

	stringifiedQuery, _ := json.Marshal(q)
	fmt.Printf("\nQuery: %v\n", string(stringifiedQuery))

	alerts, err := db.GetEnrichedAlertsByWorkflowId("",
		context.Background(),
		parsedIntegration.Inventory.AlertsCollection,
		q)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %v", err)
	}

	fmt.Printf("\n##Time range: %v - %v\n##Services: %v\n##Alerts len: %v\n",
		endTime,
		startTime,
		services,
		len(alerts))

	for _, alert := range alerts {
		output = append(output, map[string]string{
			"alertId": alert.Id.Hex(),
			"service": alert.Service,
			"state":   string(alert.State),
			"name":    alert.AlertName,
		})
	}

	return output, nil
}

func createIncident(input any, integration any) ([]any, error) {
	var parsedInput CreateIncidentInput
	var parsedAlert models.EnrichedAlert
	var severity models.IncidentSeverity
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

	// Search for assignee in db
	assignee := models.User{}

	// Filling in Tasks with Items
	tasks := make([]models.Task, 0)

	for si, step := range assertedIntegration.Inventory.WorkflowProperties.Steps {
		isDone := true

		stepOutputs, exists := parsedAlert.AdditionalContext[fmt.Sprintf("%s_%s", step.Integration, step.Name)].([]any)
		if !exists {
			isDone = false
		}
		// Check if is done
		if len(stepOutputs) == 0 {
			isDone = false
		}

		task := models.Task{
			Id:       primitive.NewObjectID(),
			Assignee: models.User{},
			IsDone:   isDone,
			Items:    make([]models.Item, 0),
			Priority: si,
			TaskName: step.DisplayName,
		}

		for _, stepOutput := range stepOutputs {
			item := models.Item{
				Content: make([]models.ItemContent, 0),
				Source:  step.Integration,
			}
			var parsedValue string
			for key, value := range stepOutput.(map[string]any) {
				_, isString := value.(string)
				if isString {
					parsedValue = value.(string)
				} else {
					valueBytes, err := json.Marshal(value)
					if err != nil {
						continue
					}
					parsedValue = string(valueBytes)
				}
				item.Content = append(item.Content, models.ItemContent{
					Key:       key,
					Value:     parsedValue,
					ValueType: "markdown",
				})
			}
			task.Items = append(task.Items, item)
		}
		tasks = append(tasks, task)
	}

	if parsedInput.Severity != "" {
		severity = models.IncidentSeverity(parsedInput.Severity)
	} else {
		severity = models.IncidentSeverityLow
	}

	incident := models.Incident{
		Id:          primitive.NewObjectID(),
		Assignee:    assignee,
		Summary:     "",
		History:     []models.IncidentUpdate[models.Update]{},
		NamespaceId: assertedIntegration.Inventory.WorkflowProperties.NamespaceId,
		Status:      models.IncidentStatusOpen,
		Severity:    severity,
		Tasks:       tasks,
		Title:       assertedIntegration.Inventory.WorkflowProperties.Name,
		Timestamp:   time.Now().Unix(),
	}

	if parsedInput.ManuallyCorrelatedAlerts != nil {
		var items = make([]models.Item, 0)
		for _, alert := range parsedInput.ManuallyCorrelatedAlerts {
			var parsedAlert models.EnrichedAlert
			err = json.Unmarshal([]byte(alert), &parsedAlert)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
			}
			alertMarkdown := fmt.Sprintf("\n*AlertId*:\n %s \n *AlertName*:\n %s \n *Timestamp*:\n %s \n *State*:\n %s",
				parsedAlert.Id,
				parsedAlert.AlertName,
				parsedAlert.StartTime,
				parsedAlert.State,
			)
			item := models.Item{
				Content: make([]models.ItemContent, 0),
				Source:  "signal0ne",
			}
			item.Content = append(item.Content, models.ItemContent{
				Key:       "alert",
				Value:     alertMarkdown,
				ValueType: "markdown",
			})
			items = append(items, item)
		}
		task := models.Task{
			Id:       primitive.NewObjectID(),
			Assignee: models.User{},
			IsDone:   true,
			Items:    items,
			Priority: len(tasks),
			TaskName: "Manually Correlated Alerts via Slack Integration",
		}
		incident.Tasks = append(incident.Tasks, task)
	}

	_, err = assertedIntegration.Inventory.IncidentCollection.InsertOne(context.Background(), incident)
	if err != nil {
		return nil, fmt.Errorf("failed to insert incident: %v", err)
	}

	cfg := config.GetInstance()
	id := incident.Id.Hex()
	//Construct metadata incident response
	output = append(output, map[string]any{
		"id":       id,
		"name":     incident.Title,
		"status":   incident.Status,
		"severity": incident.Severity,
		"url":      fmt.Sprintf("%s/incidents/%s", cfg.FrontendUrl, id),
	})

	return output, nil
}
