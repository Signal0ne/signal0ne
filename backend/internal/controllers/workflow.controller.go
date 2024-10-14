package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"signal0ne/cmd/config"
	"signal0ne/internal/db"
	"signal0ne/internal/errors"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations"
	"signal0ne/pkg/integrations/alertmanager"
	"signal0ne/pkg/integrations/backstage"
	"signal0ne/pkg/integrations/confluence"
	"signal0ne/pkg/integrations/datadog"
	"signal0ne/pkg/integrations/github"
	"signal0ne/pkg/integrations/jaeger"
	"signal0ne/pkg/integrations/openai"
	"signal0ne/pkg/integrations/opensearch"
	"signal0ne/pkg/integrations/pagerduty"
	"signal0ne/pkg/integrations/servicenow"
	"signal0ne/pkg/integrations/signal0ne"
	"signal0ne/pkg/integrations/slack"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkflowController struct {
	AlertsCollection    *mongo.Collection
	DebugLogger         *log.Logger
	IncidentsCollection *mongo.Collection
	PyInterface         net.Conn
	WebhookServerRef    config.Server
	WorkflowsCollection *mongo.Collection
	// ==== Use as readonly ====
	NamespaceCollection    *mongo.Collection
	IntegrationsCollection *mongo.Collection
	// =========================
}

func NewWorkflowController(
	workflowsCollection *mongo.Collection,
	namespaceCollection *mongo.Collection,
	integrationsCollection *mongo.Collection,
	incidentsCollection *mongo.Collection,
	alertsCollection *mongo.Collection,
	webhookServerRef config.Server,
	pyInterface net.Conn,
	debugLogger *log.Logger) *WorkflowController {
	return &WorkflowController{
		WorkflowsCollection:    workflowsCollection,
		NamespaceCollection:    namespaceCollection,
		IntegrationsCollection: integrationsCollection,
		IncidentsCollection:    incidentsCollection,
		AlertsCollection:       alertsCollection,
		WebhookServerRef:       webhookServerRef,
		PyInterface:            pyInterface,
		DebugLogger:            debugLogger,
	}
}

func (c *WorkflowController) ApplyWorkflow(ctx *gin.Context) {
	var workflow *models.Workflow
	var namespace *models.Namespace

	namespaceId := ctx.Param("namespaceid")

	err := ctx.ShouldBindJSON(&workflow)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("cannot parse body, %s", err),
		})
		return
	}

	nsID, _ := primitive.ObjectIDFromHex(namespaceId)
	res := c.NamespaceCollection.FindOne(ctx, bson.M{"_id": nsID})
	err = res.Decode(&namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("cannot find namespace, %s", err),
		})
		return
	}

	if err = c.validate(ctx, *workflow); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("validation error, err: %s", err),
		})
		return
	}

	workflow.Id = primitive.NewObjectID()
	workflow.NamespaceId = namespaceId
	workflow.WorkflowSalt, err = tools.GenerateWebhookSalt()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("cannot initialize webhook, %s", err),
		})
		return
	}

	insResult, err := c.WorkflowsCollection.InsertOne(ctx, workflow)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	workflowID, exists := insResult.InsertedID.(primitive.ObjectID)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot parse inserted id",
		})
		return
	}

	httpPrefix := "http://"
	if c.WebhookServerRef.ServerIsSecure {
		httpPrefix = "https://"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"webhook": fmt.Sprintf("%s%s:%s/api/webhook/%s/%s/%s",
			httpPrefix,
			c.WebhookServerRef.ServerDomain,
			c.WebhookServerRef.ServerPort,
			workflow.NamespaceId,
			workflowID.Hex(),
			workflow.WorkflowSalt,
		),
		"workflow": workflow,
	})
}

func (c *WorkflowController) GetWorkflow(ctx *gin.Context) {
	var namespace *models.Namespace
	var workflow models.Workflow

	namespaceId := ctx.Param("namespaceid")
	workflowId := ctx.Param("workflowid")

	nsID, _ := primitive.ObjectIDFromHex(namespaceId)
	res := c.NamespaceCollection.FindOne(ctx, bson.M{"_id": nsID})
	err := res.Decode(&namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("cannot find namespace for requested workflow, %s", err),
		})

		return
	}

	workflow, err = db.GetWorkflowById(workflowId, ctx, c.WorkflowsCollection)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "workflow not found",
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"workflow": workflow})
}

func (c *WorkflowController) GetWorkflows(ctx *gin.Context) {
	var namespace *models.Namespace
	var workflows []models.Workflow

	namespaceId := ctx.Param("namespaceid")

	nsID, _ := primitive.ObjectIDFromHex(namespaceId)
	res := c.NamespaceCollection.FindOne(ctx, bson.M{"_id": nsID})
	err := res.Decode(&namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("cannot find namespace for requested workflows, %s", err),
		})
		return
	}

	cursor, err := c.WorkflowsCollection.Find(ctx, bson.M{"namespaceId": namespaceId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("cannot find workflows, %s", err),
		})
		return
	}

	err = cursor.All(ctx, &workflows)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("cannot decode workflows, %s", err),
		})
		return
	}

	if workflows == nil {
		workflows = []models.Workflow{}
	}

	ctx.JSON(http.StatusOK, gin.H{"workflows": workflows})
}

func (c *WorkflowController) validate(ctx context.Context, workflow models.Workflow) error {

	// Lookback format
	pattern := `^\d+m$`
	lookbackRegex := regexp.MustCompile(pattern)
	matches := lookbackRegex.FindStringSubmatch(workflow.Lookback)
	if len(matches) != 1 {
		return fmt.Errorf("lookback invalid format, proper format example: '15m'")
	}

	// Trigger schema
	if workflow.Trigger.Webhook.Output != nil {
		fmt.Printf("Output Webhook Trigger: %s", workflow.Trigger.Webhook.Output)
	} else if workflow.Trigger.Scheduler.Output != nil {
		fmt.Printf("Output Scheduler Trigger: %s", workflow.Trigger.Scheduler.Output)
	} else {
		return fmt.Errorf("specified trigger type doesn't exist")
	}

	// Steps
	for index, step := range workflow.Steps {
		var integrationTemplate models.Integration
		filter := bson.M{
			"name": step.Integration,
		}
		result := c.IntegrationsCollection.FindOne(ctx, filter)
		err := result.Decode(&integrationTemplate)
		if err != nil {
			return fmt.Errorf("integration schema parsing error, %s", err)
		}

		integrationType, exists := integrations.InstallableIntegrationTypesLibrary[integrationTemplate.Type]

		workflow.Steps[index].IntegrationType = integrationTemplate.Type

		if !exists {
			return fmt.Errorf("cannot find integration type specified")
		}

		integration := reflect.New(integrationType).Elem().Interface().(models.IIntegration)

		err = integration.ValidateStep(step.Input, step.Function)
		if err != nil {
			return err
		}

	}

	return nil
}

func (c *WorkflowController) WebhookTriggerHandler(ctx *gin.Context) {
	var workflow *models.Workflow
	var executionLog = models.StepExecution{
		Outputs:  map[string]any{},
		Outcomes: []models.StepExecutionOutcome{},
	}
	var localErrorMessage = ""
	var incomingTriggerPayload map[string]any
	var alert = models.EnrichedAlert{
		Id:                primitive.NewObjectID(),
		TriggerProperties: map[string]any{},
		AdditionalContext: map[string]any{},
	}

	namespaceId := ctx.Param("namespaceid")
	workflowId := ctx.Param("workflowid")
	salt := ctx.Param("salt")

	alert.WorkflowId = workflowId

	ID, err := primitive.ObjectIDFromHex(workflowId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	filter := bson.M{"_id": ID}

	searchRes := c.WorkflowsCollection.FindOne(ctx, filter)
	err = searchRes.Decode(&workflow)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		localErrorMessage = fmt.Sprintf("%v", err)
		executionLog.Outcomes = append(executionLog.Outcomes, models.StepExecutionOutcome{
			Status:     "failure",
			LogMessage: localErrorMessage,
		})
		return
	}

	if workflow.NamespaceId != namespaceId {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		localErrorMessage = fmt.Sprintf("%v", err)
		executionLog.Outcomes = append(executionLog.Outcomes, models.StepExecutionOutcome{
			Status:     "failure",
			LogMessage: localErrorMessage,
		})
		return
	}

	if workflow.WorkflowSalt != salt {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		localErrorMessage = fmt.Sprintf("%v", err)
		executionLog.Outcomes = append(executionLog.Outcomes, models.StepExecutionOutcome{
			Status:     "failure",
			LogMessage: localErrorMessage,
		})
		return
	}

	//Trigger execution
	body, err := ctx.GetRawData()
	if err != nil || len(body) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		localErrorMessage = fmt.Sprintf("%v", err)
		executionLog.Outcomes = append(executionLog.Outcomes, models.StepExecutionOutcome{
			Status:     "failure",
			LogMessage: localErrorMessage,
		})
		tools.RecordExecution(ctx, executionLog, c.WorkflowsCollection, filter)
		return
	}

	err = json.Unmarshal(body, &incomingTriggerPayload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		localErrorMessage = fmt.Sprintf("%v", err)
		executionLog.Outcomes = append(executionLog.Outcomes, models.StepExecutionOutcome{
			Status:     "failure",
			LogMessage: localErrorMessage,
		})
		tools.RecordExecution(ctx, executionLog, c.WorkflowsCollection, filter)
		return

	}

	alert.Integration = workflow.Trigger.Webhook.Integration

	switch alert.Integration {
	case "alertmanager":
		inventory := alertmanager.NewAlertmanagerIntegrationInventory(
			c.AlertsCollection,
		)
		integration := &alertmanager.AlertmanagerIntegration{
			Inventory: inventory,
		}
		err = integration.Trigger(incomingTriggerPayload, &alert, workflow)
	case "datadog":
		inventory := datadog.NewDataDogIntegrationInventory(
			c.PyInterface,
			c.AlertsCollection,
			workflow,
		)
		integration := &datadog.DataDogIntegration{
			Inventory: inventory,
		}
		err = integration.Trigger(incomingTriggerPayload, &alert, workflow)
	}
	if err == errors.ErrConditionNotSatisfied || err == errors.ErrAlertAlreadyInactive {
		ctx.JSON(http.StatusAccepted, nil)
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		localErrorMessage = fmt.Sprintf("%v", err)
		executionLog.Outcomes = append(executionLog.Outcomes, models.StepExecutionOutcome{
			Status:     "failure",
			LogMessage: localErrorMessage,
		})
		tools.RecordExecution(ctx, executionLog, c.WorkflowsCollection, filter)
		return
	}

	// Steps execution
	for _, step := range workflow.Steps {
		var integrationTemplate models.Integration
		localErrorMessage = ""

		// 1. Get integration
		filter := bson.M{
			"name": step.Integration,
		}
		result := c.IntegrationsCollection.FindOne(ctx, filter)
		err := result.Decode(&integrationTemplate)
		if err != nil {
			localErrorMessage = fmt.Sprintf("%v", err)
			continue
		}

		_, exists := integrations.InstallableIntegrationTypesLibrary[integrationTemplate.Type]
		if !exists {
			localErrorMessage = fmt.Sprintf("%v", err)
			continue
		}

		// 2. Parse integration
		var integration any
		switch integrationTemplate.Type {
		case "alertmanager":
			integration = &alertmanager.AlertmanagerIntegration{}
		case "backstage":
			integration = &backstage.BackstageIntegration{}
		case "confluence":
			inventory := confluence.NewConfluenceIntegrationInventory(
				c.PyInterface,
			)
			integration = &confluence.ConfluenceIntegration{
				Inventory: inventory,
			}
		case "datadog":
			inventory := datadog.NewDataDogIntegrationInventory(
				c.PyInterface,
				c.AlertsCollection,
				workflow,
			)
			integration = &datadog.DataDogIntegration{
				Inventory: inventory,
			}
		case "github":
			integration = &github.GithubIntegration{}
		case "jaeger":
			inventory := jaeger.NewJaegerIntegrationInventory(
				c.PyInterface,
			)
			integration = &jaeger.JaegerIntegration{
				Inventory: inventory,
			}
		case "openai":
			integration = &openai.OpenaiIntegration{}
		case "opensearch":
			inventory := opensearch.NewOpenSearchIntegrationInventory(
				c.PyInterface,
			)
			integration = &opensearch.OpenSearchIntegration{
				Inventory: inventory,
			}
		case "pagerduty":
			inventory := pagerduty.NewPagerdutyIntegrationInventory(
				workflow,
			)
			integration = &pagerduty.PagerdutyIntegration{
				Inventory: inventory,
			}
		case "servicenow":
			integration = &servicenow.ServicenowIntegration{}
		case "signal0ne":
			inventory := signal0ne.NewSignal0neIntegrationInventory(
				c.AlertsCollection,
				c.IncidentsCollection,
				c.PyInterface,
				workflow,
			)
			integration = &signal0ne.Signal0neIntegration{
				Inventory: inventory,
			}
		case "slack":
			inventory := slack.NewSlackIntegrationInventory(workflow.Name)
			integration = &slack.SlackIntegration{
				Inventory: inventory,
			}
		default:
			integration = &models.Integration{}
		}

		err = result.Decode(integration)
		if err != nil {
			fmt.Printf("integration schema parsing error, %s", err)
			localErrorMessage = fmt.Sprintf("%v", err)
			continue
		}

		// 3. Prepare input
		for key, value := range step.Input {
			buf := new(bytes.Buffer)
			t, err := template.New("").Funcs(template.FuncMap{
				"root": func() string {
					bytes, _ := json.Marshal(alert)
					return string(bytes)
				},
				"default": func(value any, defaultValue any) any {
					if value == "" || value == nil {
						return defaultValue
					}
					return value
				},
				"date": func(timeString string, shift string, outputType string) string {
					unit := string(shift[len(shift)-1])
					value, _ := strconv.Atoi(shift[1 : len(shift)-1])
					sign := string(shift[0])
					multiplier := int64(0)
					if sign == "+" {
						multiplier = 1
					} else if sign == "-" {
						multiplier = -1
					}

					if unit == "m" {
						multiplier = multiplier * 60
					} else if unit == "h" {
						multiplier = multiplier * 3600
					}

					parsedTime, err := time.Parse(time.RFC3339, timeString)
					if err != nil {
						fmt.Printf("invalid time format: %v", err)
						localErrorMessage = fmt.Sprintf("invalid time format: %v", err)
					}

					timestamp := parsedTime.Unix()

					resultTimestamp := int64(timestamp) + int64(value)*multiplier
					resultTime := time.Unix(resultTimestamp, 0)
					if outputType == "ts" {
						return strconv.Itoa(int(resultTimestamp))
					} else if outputType == "rfc" {
						return resultTime.Format(time.RFC3339)
					}
					return strconv.Itoa(int(resultTimestamp))
				},
			}).Parse(value)
			if err != nil {
				localErrorMessage = fmt.Sprintf("%v", err)
				continue
			}
			err = t.Execute(buf, alert)
			if err != nil {
				localErrorMessage = fmt.Sprintf("%v", err)
				continue
			}
			step.Input[key] = buf.String()
			if strings.Contains(step.Input[key], "{{") {
				step.Input[key] = ""
			}
		}

		// 4. Execute
		execResult := []map[string]any{}
		if tools.EvaluateCondition(step.Condition, alert) {
			switch i := integration.(type) {
			case *backstage.BackstageIntegration:
				execResult, err = i.Execute(step.Input, step.Output, step.Function)
			case *confluence.ConfluenceIntegration:
				execResult, err = i.Execute(step.Input, step.Output, step.Function)
			case *datadog.DataDogIntegration:
				execResult, err = i.Execute(step.Input, step.Output, step.Function)
			case *github.GithubIntegration:
				execResult, err = i.Execute(step.Input, step.Output, step.Function)
			case *pagerduty.PagerdutyIntegration:
				execResult, err = i.Execute(step.Input, step.Output, step.Function)
			case *servicenow.ServicenowIntegration:
				execResult, err = i.Execute(step.Input, step.Output, step.Function)
			case *slack.SlackIntegration:
				execResult, err = i.Execute(step.Input, step.Output, step.Function)
			case *opensearch.OpenSearchIntegration:
				execResult, err = i.Execute(step.Input, step.Output, step.Function)
			case *jaeger.JaegerIntegration:
				execResult, err = i.Execute(step.Input, step.Output, step.Function)
			case *signal0ne.Signal0neIntegration:
				execResult, err = i.Execute(step.Input, step.Output, step.Function)
			case *alertmanager.AlertmanagerIntegration:
				execResult, err = i.Execute(step.Input, step.Output, step.Function)
			case *openai.OpenaiIntegration:
				execResult, err = i.Execute(step.Input, step.Output, step.Function)
			default:
				err = fmt.Errorf("unknown integration type")
			}
		}
		if err != nil {
			fmt.Printf("failed to execute %s, error: %s", step.Function, err)
			localErrorMessage = fmt.Sprintf("%v", err)

			tools.LogStep(
				c.DebugLogger,
				step.Function,
				integrationTemplate.Name,
				localErrorMessage,
				"failure",
				step.Input,
				execResult)
		} else {

			tools.LogStep(
				c.DebugLogger,
				step.Function,
				integrationTemplate.Name,
				localErrorMessage,
				"success",
				step.Input,
				execResult)
		}

		if execResult != nil {
			alert.AdditionalContext[fmt.Sprintf("%s_%s", integrationTemplate.Name, step.Name)] = execResult
			executionLog.Outputs[step.Name] = execResult
		}

		status := "success"
		if localErrorMessage != "" {
			status = "failure"
		}
		executionLog.Outcomes = append(executionLog.Outcomes, models.StepExecutionOutcome{
			Status:     status,
			LogMessage: localErrorMessage,
		})
	}

	executionLog.ParsedWorkflow = models.ParsedWorkflow{
		Steps:   workflow.Steps,
		Trigger: workflow.Trigger,
	}
	tools.RecordExecution(ctx, executionLog, c.WorkflowsCollection, filter)

	_, err = c.AlertsCollection.InsertOne(ctx, alert)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
