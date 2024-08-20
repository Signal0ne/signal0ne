package controllers

import (
	"github.com/gin-gonic/gin"
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

func (ac *AlertsController) Details(ctx *gin.Context) {}

func (ac *AlertsController) Summary(ctx *gin.Context) {}

func (ac *AlertsController) Correlations(ctx *gin.Context) {}
