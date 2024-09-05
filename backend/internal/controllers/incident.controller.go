package controllers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"signal0ne/internal/db"
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

	incident, err := db.GetIncidentById(id, ctx, ic.IncidentsCollection)
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

	integrationTemplate, err := db.GetIntegrationByName(
		createIncidentRequest.Integration,
		ctx,
		ic.IntegrationsCollection,
	)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "integration not found",
		})
		return
	}

	workflow, err := db.GetWorkflowById(createIncidentRequest.BaseAlertId,
		ctx,
		ic.WorkflowsCollection,
	)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "workflow not found",
		})
		return
	}

	alert, err := db.GetEnrichedAlertById(
		createIncidentRequest.BaseAlertId,
		ctx,
		ic.AlertsCollection,
	)
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
		db.SetRegisterHistoryEvent(
			incidentId,
			incidentUpdate,
			ctx,
			ic.IncidentsCollection,
		)
	case "task":
		var incidentUpdate models.IncidentUpdate[models.TaskUpdate]
		err := ctx.BindJSON(&incidentUpdate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("invalid request: %v", err),
			})
			return
		}
		db.SetRegisterHistoryEvent(
			incidentId,
			incidentUpdate,
			ctx,
			ic.IncidentsCollection,
		)
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
