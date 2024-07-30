package controllers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"signal0ne/internal/models"
	"signal0ne/pkg/integrations"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type IntegrationController struct {
	IntegrationCollection *mongo.Collection
	NamespaceCollection   *mongo.Collection //Must be used as Readonly
}

func NewIntegrationController(
	integrationCollection *mongo.Collection,
	namespaceCollection *mongo.Collection) *IntegrationController {
	return &IntegrationController{
		IntegrationCollection: integrationCollection,
		NamespaceCollection:   namespaceCollection,
	}
}

func (ic *IntegrationController) Install(ctx *gin.Context) {
	var integrationTemplate models.Integration
	err := ctx.ShouldBindJSON(&integrationTemplate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	integType, exists := integrations.InstallableIntegrationTypesLibrary[integrationTemplate.Type]
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot find requested integartion",
		})
	}

	body, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	integration := reflect.New(integType).Elem().Interface().(models.IIntegration)

	err = json.Unmarshal(body, integration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	_, err = ic.IntegrationCollection.InsertOne(ctx, integration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	ctx.JSON(http.StatusOK, integration)
}

func (ic *IntegrationController) ListInstalled(ctx *gin.Context) {
}

func (ic *IntegrationController) ListInstallable(ctx *gin.Context) {
	installableIntegrations, err := integrations.GetInstallableIntegrationsLib()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
	ctx.JSON(http.StatusOK, installableIntegrations)
}
