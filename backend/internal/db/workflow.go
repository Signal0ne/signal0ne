package db

import (
	"context"
	"signal0ne/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetWorkflowById(
	id string,
	ctx context.Context,
	workflowsCollection *mongo.Collection,
) (models.Workflow, error) {
	var workflow models.Workflow

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return workflow, err
	}

	err = workflowsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&workflow)
	if err != nil {
		return workflow, err
	}

	return workflow, nil
}
