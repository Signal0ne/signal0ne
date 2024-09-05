package db

import (
	"context"
	"signal0ne/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetEnrichedAlertById(
	id string,
	ctx context.Context,
	alertsCollection *mongo.Collection,
) (models.EnrichedAlert, error) {
	var alert models.EnrichedAlert

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return alert, err
	}

	err = alertsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&alert)
	if err != nil {
		return alert, err
	}

	return alert, nil
}
