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
	if err == mongo.ErrNoDocuments {
		return alert, nil
	}
	if err != nil {
		return alert, err
	}

	return alert, nil
}

func GetEnrichedAlertsByWorkflowId(
	workflowId string,
	ctx context.Context,
	alertsCollection *mongo.Collection,
	filter bson.M,
) ([]models.EnrichedAlert, error) {
	var alerts []models.EnrichedAlert

	if workflowId != "" {
		filter["workflowId"] = workflowId
	}

	cursor, err := alertsCollection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return alerts, nil
	}
	if err != nil {
		return alerts, err
	}

	err = cursor.All(ctx, &alerts)
	if err != nil {
		return alerts, err
	}

	return alerts, nil
}

func UpdateEnrichedAlert(
	alert models.EnrichedAlert,
	ctx context.Context,
	alertsCollection *mongo.Collection,
) error {
	_, err := alertsCollection.UpdateOne(
		ctx,
		bson.M{"_id": alert.Id},
		bson.M{"$set": alert},
	)
	if err != nil {
		return err
	}

	return nil
}
