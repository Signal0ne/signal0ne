package controllers

import (
	"context"
	"signal0ne/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateIncidentRequest struct {
	Integration string `json:"integration"`
	BaseAlertId string `json:"baseAlertId"`
}

type IncidentController struct {
	IncidentsCollection    *mongo.Collection
	IntegrationsCollection *mongo.Collection
}

func NewIncidentController(
	incidentsCollection *mongo.Collection,
	integrationsCollection *mongo.Collection) *IncidentController {
	return &IncidentController{
		IncidentsCollection:    incidentsCollection,
		IntegrationsCollection: integrationsCollection,
	}
}

func (ic *IncidentController) GetIncident(ctx *gin.Context) {
	id := ctx.Param("incidentid")

	incident, err := ic.getIncidentById(id, ctx)
	if err != nil {
		ctx.JSON(404, gin.H{
			"error": "incident not found",
		})
		return
	}

	ctx.JSON(200, incident)
}

func (ic *IncidentController) CreateIncident(ctx *gin.Context) {
	var createIncidentRequest CreateIncidentRequest
	err := ctx.BindJSON(&createIncidentRequest)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}

	incidentType := createIncidentRequest.Integration
	integration, err := ic.getIntegrationByIncidentType(incidentType, ctx)
	if err != nil || integration == (models.Integration{}) {
		ctx.JSON(404, gin.H{
			"error": "integration not found",
		})
		return
	}

	switch incidentType {
	case "signal0ne":
		ic.IncidentsCollection.InsertOne(context.Background(), models.Incident{})
	case "pagerduty":
		// Create incident in PagerDuty
	default: //signal0ne is also the default
		ic.IncidentsCollection.InsertOne(context.Background(), models.Incident{})
	}

}

func (ic *IncidentController) getIncidentById(id string, ctx context.Context) (models.Incident, error) {
	var incident models.Incident

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return incident, err
	}

	err = ic.IncidentsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&incident)
	if err != nil {
		return incident, err
	}

	return incident, nil
}

func (ic *IncidentController) getIntegrationByIncidentType(integrationName string, ctx context.Context) (models.Integration, error) {
	res := ic.IntegrationsCollection.FindOne(ctx, bson.M{"type": integrationName})

	var integration models.Integration
	err := res.Decode(&integration)
	if err != nil {
		return models.Integration{}, err
	}

	return integration, nil
}
