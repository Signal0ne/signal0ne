package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"signal0ne/cmd/config"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations"
	"signal0ne/pkg/integrations/alertmanager"
	"signal0ne/pkg/integrations/backstage"
	"signal0ne/pkg/integrations/jaeger"
	"signal0ne/pkg/integrations/openai"
	"signal0ne/pkg/integrations/opensearch"
	"signal0ne/pkg/integrations/signal0ne"
	"signal0ne/pkg/integrations/slack"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkflowController struct {
	WebhookServerRef    config.Server
	WorkflowsCollection *mongo.Collection
	PyInterface         net.Conn
	IncidentsCollection *mongo.Collection
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
	webhookServerRef config.Server,
	pyInterface net.Conn) *WorkflowController {
	return &WorkflowController{
		WorkflowsCollection:    workflowsCollection,
		NamespaceCollection:    namespaceCollection,
		IntegrationsCollection: integrationsCollection,
		IncidentsCollection:    incidentsCollection,
		WebhookServerRef:       webhookServerRef,
		PyInterface:            pyInterface,
	}
}

func (c *WorkflowController) WebhookTriggerHandler(ctx *gin.Context) {
	var workflow *models.Workflow
	var localErrorMessage = ""
	var alert = models.EnrichedAlert{
		Id:                primitive.NewObjectID(),
		TriggerProperties: map[string]any{},
		AdditionalContext: map[string]models.Outputs{},
	}

	namespaceId := ctx.Param("namespaceid")
	workflowId := ctx.Param("workflowid")
	salt := ctx.Param("salt")

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
		tools.RecordExecution(ctx, localErrorMessage, c.WorkflowsCollection, filter)
		return
	}

	if workflow.NamespaceId != namespaceId {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		localErrorMessage = fmt.Sprintf("%v", err)
		tools.RecordExecution(ctx, localErrorMessage, c.WorkflowsCollection, filter)
		return
	}

	if workflow.WorkflowSalt != salt {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		localErrorMessage = fmt.Sprintf("%v", err)
		tools.RecordExecution(ctx, localErrorMessage, c.WorkflowsCollection, filter)
		return
	}

	//Trigger execution
	alert.TriggerProperties, err = tools.WebhookTriggerExec(ctx, workflow)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		localErrorMessage = fmt.Sprintf("%v", err)
		tools.RecordExecution(ctx, localErrorMessage, c.WorkflowsCollection, filter)
		return
	}

	// Steps execution
	for _, step := range workflow.Steps {
		var integrationTemplate models.Integration

		// 1. Get integration
		filter := bson.M{
			"name": step.Integration,
		}
		result := c.IntegrationsCollection.FindOne(ctx, filter)
		err := result.Decode(&integrationTemplate)
		if err != nil {
			fmt.Printf("integration schema parsing error, %s", err)
			localErrorMessage = fmt.Sprintf("%v", err)
			continue
		}

		_, exists := integrations.InstallableIntegrationTypesLibrary[integrationTemplate.Type]
		if !exists {
			fmt.Printf("cannot find integration type specified")
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
		case "slack":
			inventory := slack.NewSlackIntegrationInventory(workflow.Name)
			integration = &slack.SlackIntegration{
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
		case "signal0ne":
			inventory := signal0ne.NewSignal0neIntegrationInventory(
				c.IncidentsCollection,
				c.PyInterface,
			)
			integration = &signal0ne.Signal0neIntegration{
				Inventory: inventory,
			}
		case "jaeger":
			inventory := jaeger.NewJaegerIntegrationInventory(
				c.PyInterface,
			)
			integration = &jaeger.JaegerIntegration{
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
				"index": func() string {
					bytes, _ := json.Marshal(alert)
					return string(bytes)
				},
				"default": func(value any, defaultValue any) any {
					if value == "" || value == nil {
						return defaultValue
					}
					return value
				},
				"date": func(timestamp float64, shift string, outputType string) string {
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
		}

		// 4. Execute
		execResult := []map[string]any{}
		if tools.EvaluateCondition(step.Condition, alert) {
			switch i := integration.(type) {
			case *backstage.BackstageIntegration:
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
		}

		alert.AdditionalContext[fmt.Sprintf("%s_%s", integrationTemplate.Name, step.Function)] = models.Outputs{
			Output: execResult,
		}
	}

	// New incident
	keys := make([]string, 0, len(alert.AdditionalContext))
	for k := range alert.AdditionalContext {
		keys = append(keys, k)
	}

	primaryFields := make([]map[string]any, 0)
	primaryFields = append(primaryFields, alert.TriggerProperties)

	tasks := make([]models.Task, 0)

	for si, step := range workflow.Steps {
		isDone := true
		fields := make([]models.Field, 0)
		output := alert.AdditionalContext[keys[si]].Output.([]map[string]any)
		for _, outputObject := range output {
			outputKeys := make([]string, 0, len(outputObject))
			for k := range outputObject {
				outputKeys = append(outputKeys, k)
			}
			field := models.Field{
				Key:       "",
				Source:    step.Integration,
				Value:     "",
				ValueType: "",
			}
		}
		task := models.Task{
			StepName: step.Name,
			Priority: si,
			Assignee: models.User{},
			IsDone:   isDone,
		}

		tasks = append(tasks, task)
	}

	incident := models.Incident{
		Id:            alert.Id,
		Title:         workflow.Name,
		Assignee:      models.User{},
		Severity:      "",
		PrimaryFields: primaryFields,
		Tasks:         tasks,
		History:       make([]models.IncidentUpdate[models.Update], 0),
	}

	c.IncidentsCollection.InsertOne(ctx, incident)

	tools.RecordExecution(ctx, localErrorMessage, c.WorkflowsCollection, filter)

	ctx.JSON(http.StatusOK, nil)
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

	workflowID := insResult.InsertedID.(primitive.ObjectID)

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
	})
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
	for _, step := range workflow.Steps {
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
