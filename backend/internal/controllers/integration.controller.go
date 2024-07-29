package controllers

import (
	"fmt"
	"net/http"
	"signal0ne/internal/models"
	"signal0ne/pkg/integrations"
	"signal0ne/pkg/integrations/backstage"

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
	var integration interface{}
	switch integrationTemplate.Type {
	case "backstage":
		integration = backstage.NewBackstageIntegration(integrationTemplate)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("No integration named %s found", integrationTemplate.Type),
		})
		return
	}

	_, err = ic.IntegrationCollection.InsertOne(ctx, integration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	ctx.JSON(http.StatusOK, nil)
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
