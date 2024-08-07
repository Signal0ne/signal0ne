package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"signal0ne/cmd/config"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"
	"signal0ne/pkg/integrations"
	"signal0ne/pkg/integrations/backstage"
	"signal0ne/pkg/integrations/opensearch"
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
	WebhookServerRef       config.Server
	WorkflowsCollection    *mongo.Collection
	NamespaceCollection    *mongo.Collection //Must be used as Readonly
	IntegrationsCollection *mongo.Collection //Must be used as Readonly
}

func NewWorkflowController(
	workflowsCollection *mongo.Collection,
	namespaceCollection *mongo.Collection,
	integrationsCollection *mongo.Collection,
	webhookServerRef config.Server) *WorkflowController {
	return &WorkflowController{
		WorkflowsCollection:    workflowsCollection,
		NamespaceCollection:    namespaceCollection,
		IntegrationsCollection: integrationsCollection,
		WebhookServerRef:       webhookServerRef,
	}
}

func (c *WorkflowController) WebhookTriggerHandler(ctx *gin.Context) {
	var workflow *models.Workflow
	var localErrorMessage = ""
	var alert = models.EnrichedAlert{
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
		case "backstage":
			integration = &backstage.BackstageIntegration{}
		case "slack":
			integration = &slack.SlackIntegration{}
		case "opensearch":
			integration = &opensearch.OpenSearchIntegration{}
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
		var alertEnrichmentsMap = make(map[string]any)
		for key, value := range alert.TriggerProperties {
			switch value.(type) {
			case string:
				alertEnrichmentsMap[key] = value
			case int64:
				alertEnrichmentsMap[key] = value
			case float64:
				alertEnrichmentsMap[key] = value
			default:
				bytes, err := json.Marshal(value)
				if err != nil {
					localErrorMessage = fmt.Sprintf("%v", err)
					continue
				}
				alertEnrichmentsMap[key] = string(bytes)
			}
		}
		for key, value := range alert.AdditionalContext {
			bytes, err := json.Marshal(value)
			if err != nil {
				localErrorMessage = fmt.Sprintf("%v", err)
				continue
			}
			alertEnrichmentsMap[key] = string(bytes)
		}

		for key, value := range step.Input {
			buf := new(bytes.Buffer)
			t, err := template.New("").Funcs(template.FuncMap{
				"index": func() string {
					bytes, _ := json.Marshal(alertEnrichmentsMap)
					return string(bytes)
				},
				"date": func(timestamp float64, shift string) string {
					unit := string(shift[len(shift)-1])
					value, _ := strconv.Atoi(shift[1 : len(shift)-2])
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

					return resultTime.Format(time.RFC3339)
				},
			}).Parse(value)
			if err != nil {
				localErrorMessage = fmt.Sprintf("%v", err)
				continue
			}
			err = t.Execute(buf, alertEnrichmentsMap)
			if err != nil {
				localErrorMessage = fmt.Sprintf("%v", err)
				continue
			}
			step.Input[key] = buf.String()
		}

		// 4. Execute
		execResult := []map[string]any{}
		switch i := integration.(type) {
		case *backstage.BackstageIntegration:
			execResult, err = i.Execute(step.Input, step.Output, step.Function)
		case *slack.SlackIntegration:
			execResult, err = i.Execute(step.Input, step.Output, step.Function)
		case *opensearch.OpenSearchIntegration:
			execResult, err = i.Execute(step.Input, step.Output, step.Function)
		default:
			err = fmt.Errorf("unknown integration type")
		}
		if err != nil {
			fmt.Printf("failed to execute %s, error: %s", step.Function, err)
			localErrorMessage = fmt.Sprintf("%v", err)
		}

		alert.AdditionalContext[fmt.Sprintf("%s.%s", integrationTemplate.Name, step.Function)] = models.Outputs{
			Output: execResult,
		}
	}

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
		"webhook": fmt.Sprintf("%s%s:%s/webhook/%s/%s/%s",
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
