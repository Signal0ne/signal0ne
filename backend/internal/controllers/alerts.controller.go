package controllers

import (
	"context"
	"fmt"
	"net/http"
	"signal0ne/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlertsController struct {
	AlertsCollection *mongo.Collection
}

func NewAlertsController(
	alertsCollection *mongo.Collection) *AlertsController {
	return &AlertsController{
		AlertsCollection: alertsCollection,
	}
}

func (ac *AlertsController) Details(ctx *gin.Context) {

	alertID := ctx.Param("alertid")

	alert, err := ac.getAlertById(alertID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%v", err)})
		return
	}

	results := make([]map[string]any, 0)

	for _, result := range alert.AdditionalContext["opensearch_get_log_occurrences"].Output.(primitive.A) {
		intermediateResult := make(map[string]any)
		for _, element := range result.(primitive.D) {
			intermediateResult[element.Key] = element.Value
		}
		results = append(results, intermediateResult)
	}

	ctx.JSON(http.StatusOK, results)
}

func (ac *AlertsController) Correlations(ctx *gin.Context) {

	alertID := ctx.Param("alertid")

	alert, err := ac.getAlertById(alertID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%v", err)})
		return
	}

	results := make([]map[string]any, 0)

	for _, result := range alert.AdditionalContext["alertmanager_get_relevant_alerts"].Output.(primitive.A) {
		intermediateResult := make(map[string]any)
		for _, element := range result.(primitive.D) {
			intermediateResult[element.Key] = element.Value
		}
		results = append(results, intermediateResult)
	}

	ctx.JSON(http.StatusOK, results)
}

func (ac *AlertsController) Summary(ctx *gin.Context) {}

func (ac *AlertsController) getAlertById(id string) (models.EnrichedAlert, error) {
	var alert models.EnrichedAlert
	ctx := context.Background()

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return alert, err
	}

	err = ac.AlertsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&alert)
	if err != nil {
		return alert, err
	}

	return alert, nil
}
