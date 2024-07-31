package controllers

import (
	"encoding/json"
	"fmt"
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
	var integrationTemplate map[string]interface{}

	body, err := ctx.GetRawData()
	if err != nil || len(body) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("cannot get body %s", err),
		})
		return
	}

	err = json.Unmarshal(body, &integrationTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot parse body",
		})
	}

	integType, exists := integrations.InstallableIntegrationTypesLibrary[integrationTemplate["type"].(string)]
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot find requested integartion",
		})
		return
	}

	integration := reflect.New(integType).Elem().Addr().Interface().(models.IIntegration)

	err = json.Unmarshal(body, &integration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot parse body",
		})
	}

	err = integration.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	_, err = ic.IntegrationCollection.InsertOne(ctx, integration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("cannot save intergation config: %v", err),
		})
		return
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
