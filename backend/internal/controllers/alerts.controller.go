package controllers

import (
	"context"
	"signal0ne/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IncidentController struct {
	IncidentsCollection *mongo.Collection
}

func NewIncidentController(
	incidentsCollection *mongo.Collection) *IncidentController {
	return &IncidentController{
		IncidentsCollection: incidentsCollection,
	}
}

func (ac *IncidentController) getIncidentById(id string) (models.EnrichedAlert, error) {
	var alert models.EnrichedAlert
	ctx := context.Background()

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return alert, err
	}

	err = ac.IncidentsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&alert)
	if err != nil {
		return alert, err
	}

	return alert, nil
}
