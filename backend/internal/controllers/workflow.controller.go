package controllers

import (
	"net/http"
	"signal0ne/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkflowController struct {
	WorkflowsCollection *mongo.Collection
	NamespaceCollection *mongo.Collection //Must be used as Readonly
}

func NewWorkflowController(
	workflowsCollection *mongo.Collection,
	namespaceCollection *mongo.Collection) *WorkflowController {
	return &WorkflowController{
		WorkflowsCollection: workflowsCollection,
		NamespaceCollection: namespaceCollection,
	}
}

func (c *WorkflowController) ReceiveAlert(ctx *gin.Context) {
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

	workflow.NamespaceId = namespaceId

	_, err = c.WorkflowsCollection.InsertOne(ctx, workflow)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
