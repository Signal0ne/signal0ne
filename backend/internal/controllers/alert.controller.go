package controllers

import (
	"net/http"

	"signal0ne/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlertController struct {
	AlertsCollection *mongo.Collection
}

func NewAlertController(alertsCollection *mongo.Collection) *AlertController {
	return &AlertController{
		AlertsCollection: alertsCollection,
	}
}

func (ac *AlertController) GetAlert(ctx *gin.Context) {
	var alert models.EnrichedAlert

	_ = ctx.Param("namespaceid")
	alertId := ctx.Param("alertid")

	parsedAlertId, err := primitive.ObjectIDFromHex(alertId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert id"})
		return
	}

	alertRes := ac.AlertsCollection.FindOne(ctx, bson.M{
		// "namespaceid": namespaceId,
		"_id": parsedAlertId,
	})
	if alertRes.Err() != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "alert not found"})
		return
	}

	err = alertRes.Decode(&alert)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error decoding alert"})
		return
	}

	ctx.JSON(http.StatusOK, alert)
}
