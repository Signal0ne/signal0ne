package controllers

import (
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
		return
	}

	if workflow.NamespaceId != namespaceId {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if workflow.WorkflowSalt != salt {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	//Trigger execution
	var incomingTriggerPayload map[string]any
	var desiredPropertiesWithValues = map[string]any{}

	body, err := ctx.GetRawData()
	if err != nil || len(body) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("cannot get body %s", err),
		})
		return
	}

	err = json.Unmarshal(body, &incomingTriggerPayload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("cannot decode body %s", err),
		})
		return
	}

	for k, m := range workflow.Trigger.WebhookTrigger.Webhook.Output {
		desiredPropertiesWithValues[k] = c.traverseOutput(incomingTriggerPayload, k, m)
	}

	fmt.Printf("Webhook incoming payload cheery-picked %v", desiredPropertiesWithValues)

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
		}

		integType, exists := integrations.InstallableIntegrationTypesLibrary[integrationTemplate.Type]
		if !exists {
			fmt.Printf("cannot find integration type specified")
		}

		// 2. Parse integration
		integration := reflect.New(integType).Elem().Interface().(models.IIntegration)

		// 3. execute
		_, err = integration.Execute(step.Input, step.Output, step.Function)
		if err != nil {
			fmt.Printf("failed to execute %s", step.Function)
		}
	}

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

		integType, exists := integrations.InstallableIntegrationTypesLibrary[integrationTemplate.Type]
		if !exists {
			return fmt.Errorf("cannot find integration type specified")
		}

		integration := reflect.New(integType).Elem().Interface().(models.IIntegration)

		err = integration.ValidateStep(step.Input, step.Function)
		if err != nil {
			return err
		}

	}

	return nil
}

func (c *WorkflowController) traverseOutput(
	payload any,
	desiredKey string,
	mapping string) any {

	switch v := payload.(type) {
	case map[string]any:
		for key, value := range v {
			if key == mapping {
				return value
			}
			c.traverseOutput(value, desiredKey, mapping)
		}
	case []any:
		for _, value := range v {
			c.traverseOutput(value, desiredKey, mapping)
		}
	default:
		return v
	}
	return nil
}
