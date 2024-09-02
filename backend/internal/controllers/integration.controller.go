package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"signal0ne/internal/models"
	"signal0ne/pkg/integrations"
	"sort"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (ic *IntegrationController) GetInstallableIntegrations(ctx *gin.Context) {
	var integrationsList []map[string]any

	installableIntegrations, err := integrations.GetInstallableIntegrationsLib()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	for _, integration := range installableIntegrations {

		if config, ok := integration["config"].(map[string]any); ok {
			if _, hostOk := config["host"].(string); hostOk {
				if _, portOk := config["port"].(string); portOk {
					config["url"] = "string"
					delete(config, "host")
					delete(config, "port")
				}
			}
		}

		integrationsList = append(integrationsList, integration)
	}

	if integrationsList == nil {
		integrationsList = []map[string]any{}
	}

	sort.Slice(integrationsList, func(i, j int) bool {
		return integrationsList[i]["typeName"].(string) < integrationsList[j]["typeName"].(string)
	})

	ctx.JSON(http.StatusOK, gin.H{"installableIntegrations": integrationsList})
}

func (ic *IntegrationController) GetInstalledIntegrations(ctx *gin.Context) {
	var installedIntegrations []models.Integration
	var namespace *models.Namespace

	namespaceId := ctx.Param("namespaceid")
	nsID, _ := primitive.ObjectIDFromHex(namespaceId)
	res := ic.NamespaceCollection.FindOne(ctx, primitive.M{"_id": nsID})
	err := res.Decode(&namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("cannot find namespace: %v", err),
		})
		return
	}

	cursor, err := ic.IntegrationCollection.Find(ctx, primitive.M{"namespaceId": namespaceId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("cannot get installed integrations: %v", err),
		})
		return
	}

	err = cursor.All(ctx, &installedIntegrations)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("cannot decode installed integrations: %v", err),
		})
		return
	}

	if installedIntegrations == nil {
		installedIntegrations = []models.Integration{}
	}

	ctx.JSON(http.StatusOK, gin.H{"installedIntegrations": installedIntegrations})
}

func (ic *IntegrationController) Install(ctx *gin.Context) {
	var integrationTemplate map[string]interface{}
	var namespace *models.Namespace

	namespaceId := ctx.Param("namespaceid")
	nsID, _ := primitive.ObjectIDFromHex(namespaceId)
	res := ic.NamespaceCollection.FindOne(ctx, bson.M{"_id": nsID})
	err := res.Decode(&namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("cannot find namespace, %s", err),
		})
		return
	}

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

	integrationTemplate["namespaceId"] = namespaceId

	integType, exists := integrations.InstallableIntegrationTypesLibrary[integrationTemplate["type"].(string)]
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot find requested integration",
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

	_, err = ic.IntegrationCollection.InsertOne(ctx, integrationTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("cannot save integration config: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, integration)
}
