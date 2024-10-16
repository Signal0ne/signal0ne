package controllers

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"signal0ne/internal/db"
	"signal0ne/internal/models"
	"signal0ne/pkg/integrations/signal0ne"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlertController struct {
	AlertsCollection *mongo.Collection
	PyInterface      net.Conn

	// ==== Use as readonly ====
	WorkflowsCollection *mongo.Collection
	IncidentsCollection *mongo.Collection
	// =========================
}

func NewAlertController(alertsCollection *mongo.Collection,
	incidentsCollection *mongo.Collection,
	pyInterface net.Conn,
	workflowsCollection *mongo.Collection) *AlertController {
	return &AlertController{
		AlertsCollection:    alertsCollection,
		IncidentsCollection: incidentsCollection,
		PyInterface:         pyInterface,
		WorkflowsCollection: workflowsCollection,
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

	//---------
	// The only properties we refresh on GET are outputs with tags "alerts"
	// - This is due to dynamic nature of alerts, the goal is to avoid
	// situations when alert state is out of sync with the actual state or
	// non existing alerts are correlated. If it does work well trough alfa we can
	// consider to refresh all outputs on GET
	//---------
	_ = ac.SyncCorrelateAlertsFromDiffSources(ctx, alert)

	ctx.JSON(http.StatusOK, alert)
}

func (ac *AlertController) SyncCorrelateAlertsFromDiffSources(ctx *gin.Context, alert models.EnrichedAlert) error {
	var functionKey = "correlate_ongoing_alerts"

	var workflow models.Workflow
	var dependencyMap string

	workflowId, err := primitive.ObjectIDFromHex(alert.WorkflowId)
	if err != nil {
		return err
	}

	workflowRes := ac.WorkflowsCollection.FindOne(ctx, bson.M{
		"_id": workflowId,
	})
	if workflowRes.Err() != nil {
		return workflowRes.Err()
	}

	err = workflowRes.Decode(&workflow)
	if err != nil {
		return err
	}

	for _, step := range workflow.Steps {
		if step.Function == functionKey {
			var integration *signal0ne.Signal0neIntegration
			execResult := []map[string]any{}

			ongoingAlertsOutput, _ := alert.AdditionalContext[fmt.Sprintf("%s_%s", step.Integration, step.Name)].(bson.A)
			if ongoingAlertsOutput == nil {
				return fmt.Errorf("output not found for step %s", step.Name)
			}
			dependencyMap += ongoingAlertsOutput[0].(map[string]any)["dependency_map"].(string)

			switch step.Integration {
			case "signal0ne":
				inventory := signal0ne.NewSignal0neIntegrationInventory(
					ac.AlertsCollection,
					ac.IncidentsCollection,
					ac.PyInterface,
					&workflow,
				)
				integration = &signal0ne.Signal0neIntegration{
					Inventory: inventory,
				}
				input := signal0ne.CorrelateOngoingAlertsInput{
					StartTimestamp: alert.StartTime.Add(time.Minute * 5).Format(time.RFC3339),
					DependencyMap:  dependencyMap,
				}
				execResult, err = integration.Execute(input, step.Output, step.Function)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("integration %s not found", step.Integration)
			}

			alert.AdditionalContext[fmt.Sprintf("%s_%s", step.Integration, step.Name)] = execResult

			err = db.UpdateEnrichedAlert(alert, ctx, ac.AlertsCollection)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
