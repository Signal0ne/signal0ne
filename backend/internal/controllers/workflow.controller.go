package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"signal0ne/cmd/config"
	"signal0ne/internal/models"
	"signal0ne/internal/tools"

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

func (c *WorkflowController) ReceiveAlert(ctx *gin.Context) {
	var workflow *models.Workflow

	namespaceId := ctx.Param("namespaceid")
	workflowId := ctx.Param("workflowid")
	salt := ctx.Param("salt")

	searchRes := c.WorkflowsCollection.FindOne(ctx, bson.M{"_id": workflowId})
	err := searchRes.Decode(&workflow)
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

	// TBD: Execute workflow

	ctx.JSON(http.StatusOK, nil)
}

func (c *WorkflowController) ApplyWorkflow(ctx *gin.Context) {
	var workflow *models.Workflow
	var namespace *models.Namespace

	namespaceId := ctx.Param("namespaceid")

	err := ctx.ShouldBindJSON(&workflow)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	res := c.NamespaceCollection.FindOne(ctx, bson.M{"_id": namespaceId})
	err = res.Decode(&namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err = c.validate(*workflow); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	workflow.NamespaceId = namespaceId
	workflow.WorkflowSalt, err = tools.GenerateWebhookSalt()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
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
		"webhook": fmt.Sprintf("%s%s/webhook/%s/%s/%s",
			httpPrefix,
			c.WebhookServerRef.ServerDomain,
			workflow.NamespaceId,
			workflowID.Hex(),
			workflow.WorkflowSalt,
		),
	})
}

func (c *WorkflowController) validate(workflow models.Workflow) error {

	// Lookback format
	pattern := `^(\d+)m$`
	lookbackRegex := regexp.MustCompile(pattern)
	matches := lookbackRegex.FindStringSubmatch(workflow.Lookback)
	if len(matches) != 1 {
		return fmt.Errorf("lookback invalid format, proper format example: '15m'")
	}

	// Trigger schema

	// Steps schema
	for i := 0; i < len(workflow.Steps); i++ {
		step := workflow.Steps[i]
		switch step.Integration {
		case "":
		}
	}

	return nil
}
