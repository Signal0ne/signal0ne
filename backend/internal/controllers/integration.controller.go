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

type InstalledIntegration struct {
	models.Integration `json:",inline" bson:",inline"`
	Config             map[string]any `json:"config" bson:"config"`
}

func NewIntegrationController(
	integrationCollection *mongo.Collection,
	namespaceCollection *mongo.Collection) *IntegrationController {
	return &IntegrationController{
		IntegrationCollection: integrationCollection,
		NamespaceCollection:   namespaceCollection,
	}
}

func (ic *IntegrationController) GetIntegration(ctx *gin.Context) {
	var integration InstalledIntegration
	var namespace *models.Namespace

	integrationId := ctx.Param("integrationid")
	namespaceId := ctx.Param("namespaceid")

	nsID, _ := primitive.ObjectIDFromHex(namespaceId)
	res := ic.NamespaceCollection.FindOne(ctx, primitive.M{"_id": nsID})
	err := res.Decode(&namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot find namespace: %v", err),
		})
		return
	}

	integID, _ := primitive.ObjectIDFromHex(integrationId)
	err = ic.IntegrationCollection.FindOne(ctx, primitive.M{"_id": integID, "namespaceId": namespaceId}).Decode(&integration)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot find integration: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"integration": integration})
}

func (ic *IntegrationController) GetInstallableIntegrations(ctx *gin.Context) {
	var integrationsList []map[string]any

	installableIntegrations, err := integrations.GetInstallableIntegrationsLib()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	for _, integration := range installableIntegrations {
		integrationsList = append(integrationsList, integration)
	}

	if integrationsList == nil {
		integrationsList = []map[string]any{}
	}

	// Sort integrations alphabetically by type
	sort.Slice(integrationsList, func(i, j int) bool {
		return integrationsList[i]["type"].(string) < integrationsList[j]["type"].(string)
	})

	ctx.JSON(http.StatusOK, gin.H{"installableIntegrations": integrationsList})
}

func (ic *IntegrationController) GetInstalledIntegrations(ctx *gin.Context) {
	var installedIntegrations []InstalledIntegration
	var namespace *models.Namespace

	namespaceId := ctx.Param("namespaceid")
	nsID, _ := primitive.ObjectIDFromHex(namespaceId)
	res := ic.NamespaceCollection.FindOne(ctx, primitive.M{"_id": nsID})
	err := res.Decode(&namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot find namespace: %v", err),
		})
		return
	}

	cursor, err := ic.IntegrationCollection.Find(ctx, primitive.M{"namespaceId": namespaceId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Cannot get installed integrations: %v", err),
		})
		return
	}

	err = cursor.All(ctx, &installedIntegrations)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Cannot decode installed integrations: %v", err),
		})
		return
	}

	if installedIntegrations == nil {
		installedIntegrations = []InstalledIntegration{}
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
			"error": fmt.Sprintf("Cannot find namespace, %s", err),
		})
		return
	}

	body, err := ctx.GetRawData()
	if err != nil || len(body) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot get body %s", err),
		})
		return
	}

	err = json.Unmarshal(body, &integrationTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot parse body",
		})
		return
	}

	integrationTemplate["namespaceId"] = namespaceId

	integType, exists := integrations.InstallableIntegrationTypesLibrary[integrationTemplate["type"].(string)]
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot find requested integration",
		})
		return
	}

	integration := reflect.New(integType).Elem().Addr().Interface().(models.IIntegration)

	err = json.Unmarshal(body, &integration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot parse body",
		})
		return
	}

	err = integration.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	configData := integration.Initialize()

	_, err = ic.IntegrationCollection.InsertOne(ctx, integrationTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Cannot save integration config: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"integration": integration,
		"configData":  configData,
	})
}

func (ic *IntegrationController) UpdateIntegration(ctx *gin.Context) {
	var integrationTemplate map[string]interface{}
	var namespace *models.Namespace

	integrationId := ctx.Param("integrationid")
	namespaceId := ctx.Param("namespaceid")

	nsID, _ := primitive.ObjectIDFromHex(namespaceId)
	res := ic.NamespaceCollection.FindOne(ctx, bson.M{"_id": nsID})
	err := res.Decode(&namespace)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot find namespace, %s", err),
		})
		return
	}

	body, err := ctx.GetRawData()
	if err != nil || len(body) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot get body %s", err),
		})
		return
	}

	err = json.Unmarshal(body, &integrationTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot parse body",
		})
		return
	}

	integType, exists := integrations.InstallableIntegrationTypesLibrary[integrationTemplate["type"].(string)]
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot find requested integration",
		})
		return
	}

	integration := reflect.New(integType).Elem().Addr().Interface().(models.IIntegration)

	err = json.Unmarshal(body, &integration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cannot parse body",
		})
		return
	}

	err = integration.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	configData := integration.Initialize()

	updatedIntegration := bson.M{"$set": integrationTemplate}

	integID, _ := primitive.ObjectIDFromHex(integrationId)

	_, err = ic.IntegrationCollection.UpdateOne(ctx, bson.M{"_id": integID, "namespaceId": namespaceId}, updatedIntegration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Cannot update integration config: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"integration": integration,
		"configData":  configData,
	})
}
