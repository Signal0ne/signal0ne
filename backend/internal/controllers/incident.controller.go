package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"signal0ne/internal/models"
	"signal0ne/pkg/integrations/pagerduty"
	"signal0ne/pkg/integrations/servicenow"
	"signal0ne/pkg/integrations/signal0ne"

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
	WorkflowsCollection    *mongo.Collection
	AlertsCollection       *mongo.Collection
	PyInterface            net.Conn
}

func NewIncidentController(
	incidentsCollection *mongo.Collection,
	integrationsCollection *mongo.Collection,
	alertsCollection *mongo.Collection,
	workflowsCollection *mongo.Collection,
	pyInterface net.Conn) *IncidentController {
	return &IncidentController{
		IncidentsCollection:    incidentsCollection,
		IntegrationsCollection: integrationsCollection,
		AlertsCollection:       alertsCollection,
		WorkflowsCollection:    workflowsCollection,
		PyInterface:            pyInterface,
	}
}

func (ic *IncidentController) GetIncident(ctx *gin.Context) {
	id := ctx.Param("incidentid")

	incident, err := ic.getIncidentById(id, ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "incident not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, incident)
}

func (ic *IncidentController) CreateIncident(ctx *gin.Context) {
	var createIncidentRequest CreateIncidentRequest
	err := ctx.BindJSON(&createIncidentRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	integrationTemplate, err := ic.getIntegrationByName(createIncidentRequest.Integration, ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "integration not found",
		})
		return
	}

	workflow, err := ic.getWorkflowById(createIncidentRequest.BaseAlertId, ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "workflow not found",
		})
		return
	}

	alert, err := ic.getEnrichedAlertById(createIncidentRequest.BaseAlertId, ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "alert not found",
		})
		return
	}

	var integration any

	switch integrationTemplate.Type {
	case "pagerduty":
		inventory := pagerduty.NewPagerdutyIntegrationInventory(
			&workflow,
		)
		integration = &pagerduty.PagerdutyIntegration{
			Inventory: inventory,
		}

	case "servicenow":
		// Create incident in ServiceNow

	default: //signal0ne is the default
		inventory := signal0ne.NewSignal0neIntegrationInventory(
			ic.IncidentsCollection,
			ic.PyInterface,
			&workflow,
		)
		integration = &signal0ne.Signal0neIntegration{
			Inventory: inventory,
		}
	}

	filter := bson.M{
		"name": createIncidentRequest.Integration,
	}
	result := ic.IntegrationsCollection.FindOne(ctx, filter)
	err = result.Decode(&integration)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "integration not found",
		})
	}

	switch i := integration.(type) {
	case *signal0ne.Signal0neIntegration:
		bytes, _ := json.Marshal(alert)
		input := map[string]any{
			"severity":                "",
			"assignee":                models.User{},
			"parsable_context_object": string(bytes),
		}
		i.Execute(input, nil, "create_incident")
	case *pagerduty.PagerdutyIntegration:
		bytes, _ := json.Marshal(alert)
		input := map[string]any{
			"type":                    "incident",
			"title":                   workflow.Name,
			"service_name":            alert.TriggerProperties["service"],
			"parsable_context_object": string(bytes),
		}
		i.Execute(input, nil, "create_incident")
	case *servicenow.ServicenowIntegration:
		// Create incident in ServiceNow
	}

}

func (ic *IncidentController) RegisterHistoryEvent(ctx *gin.Context) {
	incidentId := ctx.Param("incidentid")
	updateType := ctx.Param("updatetype")

	switch updateType {
	case "assignee":
		var incidentUpdate models.IncidentUpdate[models.AssigneeUpdate]
		err := ctx.BindJSON(&incidentUpdate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("invalid request: %v", err),
			})
			return
		}
		ic.registerHistoryEvent(incidentId, incidentUpdate, ctx)
	case "task":
		var incidentUpdate models.IncidentUpdate[models.TaskUpdate]
		err := ctx.BindJSON(&incidentUpdate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("invalid request: %v", err),
			})
			return
		}
		ic.registerHistoryEvent(incidentId, incidentUpdate, ctx)
	}
}

func (ic *IncidentController) UpdateIncident(ctx *gin.Context) {
	incidentId := ctx.Param("incidentid")

	var incident models.Incident
	err := ctx.BindJSON(&incident)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(incidentId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid incident ID",
		})
		return
	}

	_, err = ic.IncidentsCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": incident},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
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

func (ic *IncidentController) getWorkflowById(id string, ctx context.Context) (models.Workflow, error) {
	var workflow models.Workflow

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return workflow, err
	}

	err = ic.WorkflowsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&workflow)
	if err != nil {
		return workflow, err
	}

	return workflow, nil
}

func (ic *IncidentController) getIntegrationByName(integrationName string, ctx context.Context) (models.Integration, error) {
	res := ic.IntegrationsCollection.FindOne(ctx, bson.M{"name": integrationName})

	var integration models.Integration
	err := res.Decode(&integration)
	if err != nil {
		return models.Integration{}, err
	}

	return integration, nil
}

func (ic *IncidentController) getEnrichedAlertById(id string, ctx context.Context) (models.EnrichedAlert, error) {
	var alert models.EnrichedAlert

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return alert, err
	}

	err = ic.AlertsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&alert)
	if err != nil {
		return alert, err
	}

	return alert, nil
}

func (ic *IncidentController) registerHistoryEvent(incidentId string, incidentUpdate any, ctx *gin.Context) error {
	objectID, err := primitive.ObjectIDFromHex(incidentId)
	if err != nil {
		return err
	}
	_, err = ic.IncidentsCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$push": bson.M{"history": incidentUpdate}},
	)
	if err != nil {
		return err
	}

	return nil
}
